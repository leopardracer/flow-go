package dkg

import (
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/onflow/crypto"
	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/engine"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/flow/filter"
	"github.com/onflow/flow-go/module"
	dkgmodule "github.com/onflow/flow-go/module/dkg"
	"github.com/onflow/flow-go/state/protocol"
	"github.com/onflow/flow-go/state/protocol/events"
	"github.com/onflow/flow-go/storage"
)

// DefaultPollStep specifies the default number of views that separate two calls
// to the DKG smart-contract to read broadcast messages.
const DefaultPollStep = 10

// dkgInfo consolidates information about the current DKG protocol instance.
type dkgInfo struct {
	identities      flow.IdentitySkeletonList
	phase1FinalView uint64
	phase2FinalView uint64
	phase3FinalView uint64
	// seed must be generated for each DKG instance, using a randomness source that is independent from all other nodes.
	seed []byte
}

// ReactorEngine is an engine that reacts to chain events to start new DKG runs,
// and manage subsequent phase transitions. Any unexpected error triggers a
// panic as it would undermine the security of the protocol.
// TODO replace engine.Unit with component.Component
type ReactorEngine struct {
	events.Noop
	unit              *engine.Unit
	log               zerolog.Logger
	me                module.Local
	State             protocol.State
	dkgState          storage.DKGState
	controller        module.DKGController
	controllerFactory module.DKGControllerFactory
	viewEvents        events.Views
	pollStep          uint64
}

// NewReactorEngine return a new ReactorEngine.
func NewReactorEngine(
	log zerolog.Logger,
	me module.Local,
	state protocol.State,
	dkgState storage.DKGState,
	controllerFactory module.DKGControllerFactory,
	viewEvents events.Views,
) *ReactorEngine {

	logger := log.With().
		Str("engine", "dkg_reactor").
		Logger()

	return &ReactorEngine{
		unit:              engine.NewUnit(),
		log:               logger,
		me:                me,
		State:             state,
		dkgState:          dkgState,
		controllerFactory: controllerFactory,
		viewEvents:        viewEvents,
		pollStep:          DefaultPollStep,
	}
}

// Ready implements the module ReadyDoneAware interface. It returns a channel
// that will close when the engine has successfully started.
func (e *ReactorEngine) Ready() <-chan struct{} {
	return e.unit.Ready(func() {
		// If we are starting up in the EpochSetup phase, try to start the DKG.
		// If the DKG for this epoch has been started previously, we will exit
		// and fail this epoch's DKG.
		snap := e.State.Final()

		phase, err := snap.EpochPhase()
		if err != nil {
			// unexpected storage-level error
			// TODO use irrecoverable context
			e.log.Fatal().Err(err).Msg("failed to check epoch phase when starting DKG reactor engine")
			return
		}
		epoch, err := snap.Epochs().Current()
		if err != nil {
			// unexpected storage-level error
			// TODO use irrecoverable context
			e.log.Fatal().Err(err).Msg("failed to retrieve current epoch when starting DKG reactor engine")
			return
		}
		currentCounter := epoch.Counter()
		first, err := snap.Head()
		if err != nil {
			// unexpected storage-level error
			// TODO use irrecoverable context
			e.log.Fatal().Err(err).Msg("failed to retrieve finalized header when starting DKG reactor engine")
			return
		}

		// If we start up in EpochSetup phase, attempt to start the DKG in case it wasn't started previously
		if phase == flow.EpochPhaseSetup {
			e.startDKGForEpoch(currentCounter, first)
		} else if phase == flow.EpochPhaseCommitted {
			// If we start up in EpochCommitted phase, ensure the DKG current state is set correctly.
			e.handleEpochCommittedPhaseStarted(currentCounter, first)
		}
	})
}

// Done implements the module ReadyDoneAware interface. It returns a channel
// that will close when the engine has successfully stopped.
func (e *ReactorEngine) Done() <-chan struct{} {
	return e.unit.Done()
}

// EpochSetupPhaseStarted handles the EpochSetupPhaseStarted protocol event by
// starting the DKG process.
// NOTE: ReactorEngine will not recover from mid-DKG crashes, therefore we do not need to handle dropped protocol events here.
func (e *ReactorEngine) EpochSetupPhaseStarted(currentEpochCounter uint64, first *flow.Header) {
	e.startDKGForEpoch(currentEpochCounter, first)
}

// EpochCommittedPhaseStarted handles the EpochCommittedPhaseStarted protocol
// event by checking the consistency of our locally computed key share.
// NOTE: ReactorEngine will not recover from mid-DKG crashes, therefore we do not need to handle dropped protocol events here.
func (e *ReactorEngine) EpochCommittedPhaseStarted(currentEpochCounter uint64, first *flow.Header) {
	e.handleEpochCommittedPhaseStarted(currentEpochCounter, first)
}

// startDKGForEpoch attempts to start the DKG instance for the given epoch,
// only if we have never started the DKG during setup phase for the given epoch.
// This allows consensus nodes which boot from a state snapshot within the
// EpochSetup phase to run the DKG.
//
// It starts a new controller for the epoch and registers the triggers to regularly
// query the DKG smart-contract and transition between phases at the specified views.
func (e *ReactorEngine) startDKGForEpoch(currentEpochCounter uint64, first *flow.Header) {

	firstID := first.ID()
	nextEpochCounter := currentEpochCounter + 1
	log := e.log.With().
		Uint64("cur_epoch", currentEpochCounter). // the epoch we are in the middle of
		Uint64("next_epoch", nextEpochCounter).   // the epoch we are running the DKG for
		Uint64("first_block_view", first.View).   // view of first block in EpochSetup phase
		Hex("first_block_id", firstID[:]).        // id of first block in EpochSetup phase
		Logger()

	// if we have started the dkg for this epoch already, exit
	started, err := e.dkgState.IsDKGStarted(nextEpochCounter)
	if err != nil {
		// unexpected storage-level error
		// TODO use irrecoverable context
		log.Fatal().Err(err).Msg("could not check whether DKG is started")
	}
	if started {
		log.Warn().Msg("DKG started before, skipping starting the DKG for this epoch")
		return
	}

	// flag that we are starting the dkg for this epoch
	err = e.dkgState.SetDKGState(nextEpochCounter, flow.DKGStateStarted)
	if err != nil {
		// unexpected storage-level error
		// TODO use irrecoverable context
		log.Fatal().Err(err).Msg("could not transition DKG state machine into state DKGStateStarted")
	}

	curDKGInfo, err := e.getDKGInfo(firstID)
	if err != nil {
		// unexpected storage-level error
		// TODO use irrecoverable context
		log.Fatal().Err(err).Msg("could not retrieve epoch info")
	}

	committee := curDKGInfo.identities.Filter(filter.IsConsensusCommitteeMember)

	log.Info().
		Uint64("phase1", curDKGInfo.phase1FinalView).
		Uint64("phase2", curDKGInfo.phase2FinalView).
		Uint64("phase3", curDKGInfo.phase3FinalView).
		Interface("members", committee.NodeIDs()).
		Msg("epoch info")

	if _, ok := committee.GetIndex(e.me.NodeID()); !ok {
		// node not found in DKG committee bypass starting the DKG
		log.Warn().Str("node_id", e.me.NodeID().String()).Msg("failed to find our node ID in the DKG committee skip starting DKG engine, this node will not participate in consensus after the next epoch starts")
		return
	}
	controller, err := e.controllerFactory.Create(
		dkgmodule.CanonicalInstanceID(first.ChainID, nextEpochCounter),
		committee,
		curDKGInfo.seed,
	)
	if err != nil {
		// no expected errors in controller factory
		// TODO use irrecoverable context
		log.Fatal().Err(err).Msg("could not create DKG controller")
	}
	e.controller = controller

	e.unit.Launch(func() {
		log.Info().Msg("DKG Run")
		err := e.controller.Run()
		if err != nil {
			// TODO handle crypto sentinels and do not crash here
			log.Fatal().Err(err).Msg("DKG Run error")
		}
	})

	// NOTE:
	// We register two callbacks for views that mark a state transition: one for
	// polling broadcast messages, and one for triggering the phase transition.
	// It is essential that all polled broadcast messages are processed before
	// starting the phase transition. Here we register the polling callback
	// before the phase transition, which guarantees that it will be called
	// before because callbacks for the same views are executed on a FIFO basis.
	// Moreover, the poll callback does not return until all received messages
	// are processed by the underlying DKG controller (as guaranteed by the
	// specifications and implementations of the DKGBroker and DKGController
	// interfaces).

	for view := curDKGInfo.phase1FinalView; view > first.View; view -= e.pollStep {
		e.registerPoll(view)
	}
	e.registerPhaseTransition(curDKGInfo.phase1FinalView, dkgmodule.Phase1, e.controller.EndPhase1)

	for view := curDKGInfo.phase2FinalView; view > curDKGInfo.phase1FinalView; view -= e.pollStep {
		e.registerPoll(view)
	}
	e.registerPhaseTransition(curDKGInfo.phase2FinalView, dkgmodule.Phase2, e.controller.EndPhase2)

	for view := curDKGInfo.phase3FinalView; view > curDKGInfo.phase2FinalView; view -= e.pollStep {
		e.registerPoll(view)
	}
	e.registerPhaseTransition(curDKGInfo.phase3FinalView, dkgmodule.Phase3, e.end(nextEpochCounter))
}

// handleEpochCommittedPhaseStarted is invoked upon the transition to the EpochCommitted
// phase, when the canonical beacon key vector is incorporated into the protocol state.
// Alternatively we invoke this function preemptively on startup if we are in the
// EpochCommitted Phase, in case the `EpochCommittedPhaseStarted` event was missed
// due to a crash.
//
// This function checks that the local DKG completed and that our locally computed
// key share is consistent with the canonical key vector. When this function returns,
// the current state for the just-completed DKG is guaranteed to be stored (if not, the
// program will crash). Since this function is invoked synchronously before the end
// of the current epoch, this guarantees that when we reach the end of the current epoch
// we will either have a usable beacon key committed (state [flow.RandomBeaconKeyCommitted])
// or we persist [flow.DKGStateFailure], so we can safely fall back to using our staking key.
//
// CAUTION: This function is not safe for concurrent use. This is not enforced within
// the ReactorEngine - instead we rely on the protocol event emission being single-threaded
func (e *ReactorEngine) handleEpochCommittedPhaseStarted(currentEpochCounter uint64, firstBlock *flow.Header) {

	// the DKG we have just completed produces keys that we will use in the next epoch
	nextEpochCounter := currentEpochCounter + 1

	log := e.log.With().
		Uint64("cur_epoch", currentEpochCounter). // the epoch we are in the middle of
		Uint64("next_epoch", nextEpochCounter).   // the epoch the just-finished DKG was preparing for
		Logger()

	// Check whether we have already set the current state for this DKG.
	// This can happen if the DKG failed locally, if we failed to generate
	// a local private beacon key, or if we crashed while performing this
	// check previously.
	currentState, err := e.dkgState.GetDKGState(nextEpochCounter)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Warn().Msg("failed to get dkg state, assuming this node has skipped epoch setup phase")
		} else {
			log.Fatal().Err(err).Msg("failed to get dkg state")
		}

		return
	}
	// (i) if I have a key (currentState == flow.DKGStateCompleted) which is consistent with the EpochCommit service event,
	// then commit the key or (ii) if the key is already committed (currentState == flow.RandomBeaconKeyCommitted), then we
	// expect it to be consistent with the EpochCommit service event. While (ii) is a sanity check, we have a severe problem
	// if it is violated, because a node signing with an invalid Random beacon key will be slashed - so we better check!
	// Our logic for committing a key is idempotent: it is a no-op when stating that a key `k` should be committed that previously
	// has already been committed; while it errors if `k` is different from the previously-committed key. In other words, the
	// sanity check (ii) is already included in the happy-path logic for (i). So we just repeat the happy-path logic also for
	// currentState == flow.RandomBeaconKeyCommitted, because repeated calls only occur due to node crashes, which are rare.
	if currentState != flow.DKGStateCompleted && currentState != flow.RandomBeaconKeyCommitted {
		log.Warn().Msgf("checking beacon key consistency after EpochCommit: exiting because dkg didn't reach completed state: %s", currentState.String())
		return
	}

	// Since epoch phase transitions are emitted when the first block of the new
	// phase is finalized, the block's snapshot is guaranteed to already be
	// accessible in the protocol state at this point
	snapshot := e.State.AtBlockID(firstBlock.ID())
	nextEpoch, err := snapshot.Epochs().NextCommitted()
	if err != nil {
		// CAUTION: this should never happen, indicates a storage failure or state corruption
		// TODO use irrecoverable context
		log.Fatal().Err(err).Msg("checking beacon key consistency: could not get next committed epoch")
	}
	nextDKG, err := nextEpoch.DKG()
	if err != nil {
		// CAUTION: this should never happen, indicates a storage failure or state corruption
		// TODO use irrecoverable context
		log.Fatal().Err(err).Msg("checking beacon key consistency: could not retrieve next DKG info")
		return
	}

	myBeaconPrivKey, err := e.dkgState.UnsafeRetrieveMyBeaconPrivateKey(nextEpochCounter)
	if errors.Is(err, storage.ErrNotFound) {
		log.Warn().Msg("checking beacon key consistency: no key found")
		err := e.dkgState.SetDKGState(nextEpochCounter, flow.DKGStateFailure)
		if err != nil {
			// TODO use irrecoverable context
			log.Fatal().Err(err).Msg("failed to set dkg state")
		}
		return
	} else if err != nil {
		// TODO use irrecoverable context
		log.Fatal().Err(err).Msg("checking beacon key consistency: could not retrieve beacon private key for next epoch")
		return
	}

	nextDKGPubKey, err := nextDKG.KeyShare(e.me.NodeID())
	if err != nil {
		// TODO use irrecoverable context
		log.Fatal().Err(err).Msg("checking beacon key consistency: could not retrieve my beacon public key for next epoch")
		return
	}
	localPubKey := myBeaconPrivKey.PublicKey()

	// we computed a local beacon key, but it is inconsistent with our canonical
	// public key - therefore it is unsafe for use
	if !nextDKGPubKey.Equals(localPubKey) {
		log.Warn().
			Str("computed_beacon_pub_key", localPubKey.String()).
			Str("canonical_beacon_pub_key", nextDKGPubKey.String()).
			Msg("checking beacon key consistency: locally computed beacon public key does not match beacon public key for next epoch")
		err := e.dkgState.SetDKGState(nextEpochCounter, flow.DKGStateFailure)
		if err != nil {
			// TODO use irrecoverable context
			log.Fatal().Err(err).Msg("failed to set dkg current state")
		}
		return
	}

	epochProtocolState, err := snapshot.EpochProtocolState()
	if err != nil {
		// TODO use irrecoverable context
		log.Fatal().Err(err).Msg("failed to retrieve epoch protocol state")
		return
	}
	err = e.dkgState.CommitMyBeaconPrivateKey(nextEpochCounter, epochProtocolState.Entry().NextEpochCommit)
	if err != nil {
		// TODO use irrecoverable context
		e.log.Fatal().Err(err).Msg("failed to set dkg current state")
	}
	log.Info().Msgf("successfully ended DKG, my beacon pub key for epoch %d is %s", nextEpochCounter, localPubKey)
}

// getDKGInfo returns the information required to initiate the DKG for the current epoch.
// firstBlockID must be the first block of the EpochSetup phase. This is one of the few places
// where we have to use the configuration for a future epoch that has not yet been committed.
// CAUTION: the epoch transition might not happen as described here!
// No errors are expected during normal operation.
func (e *ReactorEngine) getDKGInfo(firstBlockID flow.Identifier) (*dkgInfo, error) {
	epochsAtBlock := e.State.AtBlockID(firstBlockID).Epochs()
	currEpoch, err := epochsAtBlock.Current()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve current epoch: %w", err)
	}
	nextEpoch, err := epochsAtBlock.NextUnsafe()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve next epoch: %w", err)
	}

	seed := make([]byte, crypto.KeyGenSeedMinLen)
	_, err = rand.Read(seed)
	if err != nil {
		return nil, fmt.Errorf("could not generate random seed: %w", err)
	}

	info := &dkgInfo{
		identities:      nextEpoch.InitialIdentities(),
		phase1FinalView: currEpoch.DKGPhase1FinalView(),
		phase2FinalView: currEpoch.DKGPhase2FinalView(),
		phase3FinalView: currEpoch.DKGPhase3FinalView(),
		seed:            seed,
	}
	return info, nil
}

// registerPoll instructs the engine to query the DKG smart-contract for new
// broadcast messages at the specified view.
func (e *ReactorEngine) registerPoll(view uint64) {
	e.viewEvents.OnView(view, func(header *flow.Header) {
		e.unit.Launch(func() {
			e.unit.Lock()
			defer e.unit.Unlock()

			blockID := header.ID()
			log := e.log.With().
				Uint64("view", view).
				Uint64("height", header.Height).
				Hex("block_id", blockID[:]).
				Logger()

			log.Info().Msg("polling DKG smart-contract...")
			err := e.controller.Poll(header.ID())
			if err != nil {
				log.Err(err).Msg("failed to poll DKG smart-contract")
			}
		})
	})
}

// registerPhaseTransition instructs the engine to change phases at the
// specified view.
func (e *ReactorEngine) registerPhaseTransition(view uint64, fromState dkgmodule.State, phaseTransition func() error) {
	e.viewEvents.OnView(view, func(header *flow.Header) {
		e.unit.Launch(func() {
			e.unit.Lock()
			defer e.unit.Unlock()

			blockID := header.ID()
			log := e.log.With().
				Uint64("view", view).
				Hex("block_id", blockID[:]).
				Logger()

			log.Info().Msgf("ending %s...", fromState)
			err := phaseTransition()
			if err != nil {
				// TODO use irrecoverable context
				log.Fatal().Err(err).Msgf("node failed to end %s", fromState)
			}
			log.Info().Msgf("ended %s successfully", fromState)
		})
	})
}

// end returns a callback that is used to end the DKG protocol, save the
// resulting private key to storage, and publish the other results to the DKG
// smart-contract.
func (e *ReactorEngine) end(nextEpochCounter uint64) func() error {
	return func() error {

		err := e.controller.End()
		if crypto.IsDKGFailureError(err) {
			// Failing to complete the DKG protocol is a rare but expected scenario, which we must handle.
			// By convention, if we are leaving the happy path, we want to persist the _first_ failure symptom
			// in the `dkgState`. If the write yields a [storage.InvalidDKGStateTransitionError], it means that the state machine
			// is in the terminal state([flow.RandomBeaconKeyCommitted]) as all other transitions(even to [flow.DKGStateFailure] -> [flow.DKGStateFailure])
			// are allowed. If the protocol is in terminal state, and we have a failure symptom, then it means that recovery has happened
			// before ending the DKG. In this case, we want to ignore the error and return without error.
			e.log.Warn().Err(err).Msgf("node %s with index %d failed DKG locally", e.me.NodeID(), e.controller.GetIndex())
			err := e.dkgState.SetDKGState(nextEpochCounter, flow.DKGStateFailure)
			if err != nil {
				if storage.IsInvalidDKGStateTransitionError(err) {
					return nil
				}
				return fmt.Errorf("failed to set dkg current state following dkg end error: %w", err)
			}
			return nil // local DKG protocol has failed (the expected scenario)
		} else if err != nil {
			return fmt.Errorf("unknown error ending the dkg: %w", err)
		}

		// The following only implements the happy path, which is an atomic step-by-step progression
		// along a single path in the `dkgState` machine. If the write yields a `storage.ErrAlreadyExists`,
		// we know the overall protocol has already abandoned the happy path, because on the happy path
		// ReactorEngine is the only writer. Then this function just stops and returns without error.
		privateShare, _, _ := e.controller.GetArtifacts()
		if privateShare != nil {
			// we only store our key if one was computed
			err = e.dkgState.InsertMyBeaconPrivateKey(nextEpochCounter, privateShare)
			if err != nil {
				if errors.Is(err, storage.ErrAlreadyExists) {
					return nil // the beacon key already existing is expected in case of epoch recovery
				}
				return fmt.Errorf("could not save beacon private key in db: %w", err)
			}
		}

		err = e.controller.SubmitResult()
		if err != nil {
			return fmt.Errorf("couldn't publish DKG results: %w", err)
		}

		return nil
	}
}
