package protocol_state

import (
	"fmt"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/state/protocol"
	"github.com/onflow/flow-go/state/protocol/inmem"
)

type initialProtocolStateAdapter struct {
	*flow.RichProtocolStateEntry
	clustering flow.ClusterList
	dkg        protocol.DKG
}

var _ protocol.InitialProtocolState = (*initialProtocolStateAdapter)(nil)

func newInitialProtocolStateAdaptor(entry *flow.RichProtocolStateEntry) (*initialProtocolStateAdapter, error) {
	dkg, err := inmem.EncodableDKGFromEvents(entry.CurrentEpochSetup, entry.CurrentEpochCommit)
	if err != nil {
		return nil, fmt.Errorf("could not construct encodable DKG from events: %w", err)
	}

	clustering, err := inmem.ClusteringFromSetupEvent(entry.CurrentEpochSetup)
	if err != nil {
		return nil, fmt.Errorf("could not extract cluster list from setup event: %w", err)
	}

	return &initialProtocolStateAdapter{
		RichProtocolStateEntry: entry,
		clustering:             clustering,
		dkg:                    inmem.NewDKG(dkg),
	}, nil
}

func (s *initialProtocolStateAdapter) Epoch() uint64 {
	return s.CurrentEpochSetup.Counter
}

func (s *initialProtocolStateAdapter) Clustering() flow.ClusterList {
	return s.clustering
}

func (s *initialProtocolStateAdapter) EpochSetup() *flow.EpochSetup {
	return s.CurrentEpochSetup
}

func (s *initialProtocolStateAdapter) EpochCommit() *flow.EpochCommit {
	return s.CurrentEpochCommit
}

func (s *initialProtocolStateAdapter) DKG() protocol.DKG {
	return s.dkg
}
