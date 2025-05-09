package query

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/onflow/flow-go/fvm/errors"

	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/fvm"
	"github.com/onflow/flow-go/fvm/storage/derived"
	"github.com/onflow/flow-go/fvm/storage/snapshot"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/state/protocol"
	"github.com/onflow/flow-go/utils/debug"
	"github.com/onflow/flow-go/utils/rand"
)

const (
	DefaultLogTimeThreshold    = 1 * time.Second
	DefaultExecutionTimeLimit  = 10 * time.Second
	DefaultMaxErrorMessageSize = 1000 // 1000 chars
)

type Executor interface {
	ExecuteScript(
		ctx context.Context,
		script []byte,
		arguments [][]byte,
		blockHeader *flow.Header,
		snapshot snapshot.StorageSnapshot,
	) (
		[]byte,
		uint64,
		error,
	)

	GetAccount(
		ctx context.Context,
		addr flow.Address,
		header *flow.Header,
		snapshot snapshot.StorageSnapshot,
	) (
		*flow.Account,
		error,
	)

	GetAccountBalance(
		ctx context.Context,
		addr flow.Address,
		header *flow.Header,
		snapshot snapshot.StorageSnapshot,
	) (
		uint64,
		error,
	)

	GetAccountAvailableBalance(
		ctx context.Context,
		addr flow.Address,
		header *flow.Header,
		snapshot snapshot.StorageSnapshot,
	) (
		uint64,
		error,
	)

	GetAccountKeys(
		ctx context.Context,
		addr flow.Address,
		header *flow.Header,
		snapshot snapshot.StorageSnapshot,
	) (
		[]flow.AccountPublicKey,
		error,
	)

	GetAccountKey(
		ctx context.Context,
		addr flow.Address,
		keyIndex uint32,
		header *flow.Header,
		snapshot snapshot.StorageSnapshot,
	) (
		*flow.AccountPublicKey,
		error,
	)
}

type QueryConfig struct {
	LogTimeThreshold    time.Duration
	ExecutionTimeLimit  time.Duration
	ComputationLimit    uint64
	MaxErrorMessageSize int
}

func NewDefaultConfig() QueryConfig {
	return QueryConfig{
		LogTimeThreshold:    DefaultLogTimeThreshold,
		ExecutionTimeLimit:  DefaultExecutionTimeLimit,
		ComputationLimit:    fvm.DefaultComputationLimit,
		MaxErrorMessageSize: DefaultMaxErrorMessageSize,
	}
}

type QueryExecutor struct {
	config                QueryConfig
	logger                zerolog.Logger
	metrics               module.ExecutionMetrics
	vm                    fvm.VM
	vmCtx                 fvm.Context
	derivedChainData      *derived.DerivedChainData
	rngLock               *sync.Mutex
	protocolStateSnapshot protocol.SnapshotExecutionSubsetProvider
}

var _ Executor = &QueryExecutor{}

func NewQueryExecutor(
	config QueryConfig,
	logger zerolog.Logger,
	metrics module.ExecutionMetrics,
	vm fvm.VM,
	vmCtx fvm.Context,
	derivedChainData *derived.DerivedChainData,
	protocolStateSnapshot protocol.SnapshotExecutionSubsetProvider,
) *QueryExecutor {
	if config.ComputationLimit > 0 {
		vmCtx = fvm.NewContextFromParent(vmCtx, fvm.WithComputationLimit(config.ComputationLimit))
	}
	return &QueryExecutor{
		config:                config,
		logger:                logger,
		metrics:               metrics,
		vm:                    vm,
		vmCtx:                 vmCtx,
		derivedChainData:      derivedChainData,
		rngLock:               &sync.Mutex{},
		protocolStateSnapshot: protocolStateSnapshot,
	}
}

func (e *QueryExecutor) ExecuteScript(
	ctx context.Context,
	script []byte,
	arguments [][]byte,
	blockHeader *flow.Header,
	snapshot snapshot.StorageSnapshot,
) (
	encodedValue []byte,
	computationUsed uint64,
	err error,
) {

	startedAt := time.Now()
	memAllocBefore := debug.GetHeapAllocsBytes()

	// allocate a random ID to be able to track this script when its done,
	// scripts might not be unique so we use this extra tracker to follow their logs
	// TODO: this is a temporary measure, we could remove this in the future
	if e.logger.Debug().Enabled() {
		e.rngLock.Lock()
		defer e.rngLock.Unlock()
		trackerID, err := rand.Uint32()
		if err != nil {
			return nil, 0, fmt.Errorf("failed to generate trackerID: %w", err)
		}

		trackedLogger := e.logger.With().Hex("script_hex", script).Uint32("trackerID", trackerID).Logger()
		trackedLogger.Debug().Msg("script is sent for execution")
		defer func() {
			trackedLogger.Debug().Msg("script execution is complete")
		}()
	}

	requestCtx, cancel := context.WithTimeout(ctx, e.config.ExecutionTimeLimit)
	defer cancel()

	defer func() {
		prepareLog := func() *zerolog.Event {
			args := make([]string, 0, len(arguments))
			for _, a := range arguments {
				args = append(args, hex.EncodeToString(a))
			}
			return e.logger.Error().
				Hex("script_hex", script).
				Str("args", strings.Join(args, ","))
		}

		elapsed := time.Since(startedAt)

		if r := recover(); r != nil {
			prepareLog().
				Interface("recovered", r).
				Msg("script execution caused runtime panic")

			err = fmt.Errorf("cadence runtime error: %s", r)
			return
		}
		if elapsed >= e.config.LogTimeThreshold {
			prepareLog().
				Dur("duration", elapsed).
				Msg("script execution exceeded threshold")
		}
	}()

	var output fvm.ProcedureOutput
	_, output, err = e.vm.Run(
		fvm.NewContextFromParent(
			e.vmCtx,
			fvm.WithBlockHeader(blockHeader),
			fvm.WithProtocolStateSnapshot(e.protocolStateSnapshot.AtBlockID(blockHeader.ID())),
			fvm.WithDerivedBlockData(
				e.derivedChainData.NewDerivedBlockDataForScript(blockHeader.ID()))),
		fvm.NewScriptWithContextAndArgs(script, requestCtx, arguments...),
		snapshot)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute script (internal error): %w", err)
	}

	if output.Err != nil {
		return nil, 0, errors.NewCodedError(
			output.Err.Code(),
			"failed to execute script at block (%s): %s", blockHeader.ID(),
			summarizeLog(output.Err.Error(), e.config.MaxErrorMessageSize),
		)
	}

	encodedValue, err = jsoncdc.Encode(output.Value)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to encode runtime value: %w", err)
	}

	memAllocAfter := debug.GetHeapAllocsBytes()
	e.metrics.ExecutionScriptExecuted(
		time.Since(startedAt),
		output.ComputationUsed,
		memAllocAfter-memAllocBefore,
		output.MemoryEstimate)

	return encodedValue, output.ComputationUsed, nil
}

func summarizeLog(log string, limit int) string {
	if limit > 0 && len(log) > limit {
		split := int(limit/2) - 1
		var sb strings.Builder
		sb.WriteString(log[:split])
		sb.WriteString(" ... ")
		sb.WriteString(log[len(log)-split:])
		return sb.String()
	}
	return log
}

func (e *QueryExecutor) GetAccount(
	_ context.Context,
	address flow.Address,
	blockHeader *flow.Header,
	snapshot snapshot.StorageSnapshot,
) (
	*flow.Account,
	error,
) {
	// TODO(ramtin): utilize ctx
	blockCtx := fvm.NewContextFromParent(
		e.vmCtx,
		fvm.WithBlockHeader(blockHeader),
		fvm.WithDerivedBlockData(
			e.derivedChainData.NewDerivedBlockDataForScript(blockHeader.ID())))

	account, err := fvm.GetAccount(
		blockCtx,
		address,
		snapshot)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get account (%s) at block (%s): %w",
			address.String(),
			blockHeader.ID(),
			err)
	}

	return account, nil
}

func (e *QueryExecutor) GetAccountBalance(
	_ context.Context,
	address flow.Address,
	blockHeader *flow.Header,
	snapshot snapshot.StorageSnapshot,
) (uint64, error) {

	// TODO(ramtin): utilize ctx
	blockCtx := fvm.NewContextFromParent(
		e.vmCtx,
		fvm.WithBlockHeader(blockHeader),
		fvm.WithDerivedBlockData(
			e.derivedChainData.NewDerivedBlockDataForScript(blockHeader.ID())))

	accountBalance, err := fvm.GetAccountBalance(
		blockCtx,
		address,
		snapshot)

	if err != nil {
		return 0, fmt.Errorf(
			"failed to get account balance (%s) at block (%s): %w",
			address.String(),
			blockHeader.ID(),
			err)
	}

	return accountBalance, nil
}

func (e *QueryExecutor) GetAccountAvailableBalance(
	_ context.Context,
	address flow.Address,
	blockHeader *flow.Header,
	snapshot snapshot.StorageSnapshot,
) (uint64, error) {

	// TODO(ramtin): utilize ctx
	blockCtx := fvm.NewContextFromParent(
		e.vmCtx,
		fvm.WithBlockHeader(blockHeader),
		fvm.WithDerivedBlockData(
			e.derivedChainData.NewDerivedBlockDataForScript(blockHeader.ID())))

	accountAvailableBalance, err := fvm.GetAccountAvailableBalance(
		blockCtx,
		address,
		snapshot)

	if err != nil {
		return 0, fmt.Errorf(
			"failed to get account available balance (%s) at block (%s): %w",
			address.String(),
			blockHeader.ID(),
			err)
	}

	return accountAvailableBalance, nil
}

func (e *QueryExecutor) GetAccountKeys(
	_ context.Context,
	address flow.Address,
	blockHeader *flow.Header,
	snapshot snapshot.StorageSnapshot,
) ([]flow.AccountPublicKey, error) {
	// TODO(ramtin): utilize ctx
	blockCtx := fvm.NewContextFromParent(
		e.vmCtx,
		fvm.WithBlockHeader(blockHeader),
		fvm.WithDerivedBlockData(
			e.derivedChainData.NewDerivedBlockDataForScript(blockHeader.ID())))

	accountKeys, err := fvm.GetAccountKeys(blockCtx,
		address,
		snapshot)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get account keys (%s) at block (%s): %w",
			address.String(),
			blockHeader.ID(),
			err)
	}

	return accountKeys, nil
}

func (e *QueryExecutor) GetAccountKey(
	_ context.Context,
	address flow.Address,
	keyIndex uint32,
	blockHeader *flow.Header,
	snapshot snapshot.StorageSnapshot,
) (*flow.AccountPublicKey, error) {
	// TODO(ramtin): utilize ctx
	blockCtx := fvm.NewContextFromParent(
		e.vmCtx,
		fvm.WithBlockHeader(blockHeader),
		fvm.WithDerivedBlockData(
			e.derivedChainData.NewDerivedBlockDataForScript(blockHeader.ID())))

	accountKey, err := fvm.GetAccountKey(blockCtx,
		address,
		keyIndex,
		snapshot)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get account key (%s) at block (%s): %w",
			address.String(),
			blockHeader.ID(),
			err)
	}

	return accountKey, nil
}
