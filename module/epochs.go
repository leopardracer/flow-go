package module

import (
	"context"

	"github.com/onflow/flow-go/consensus/hotstuff/model"
	"github.com/onflow/flow-go/state/protocol"
)

// ClusterRootQCVoter is responsible for submitting a vote to the cluster QC
// contract to coordinate generation of a valid root quorum certificate for the
// next epoch.
type ClusterRootQCVoter interface {

	// Vote handles the full procedure of generating a vote, submitting it to the epoch
	// smart contract, and verifying submission. This logic should run as part of the Epoch
	// Setup Phase, i.e. at a time when the next epoch has not yet been committed. Hence,
	// this function takes the [protocol.TentativeEpoch] information as input.
	// It is safe to run `Vote` multiple times within a single Epoch Setup Phase.
	// CAUTION: epoch transition might not happen as described by [protocol.TentativeEpoch].
	// Error returns:
	//   - epochs.ClusterQCNoVoteError if we fail to vote for a benign reason
	//   - generic error in case of critical unexpected failure
	Vote(context.Context, protocol.TentativeEpoch) error
}

// QCContractClient enables interacting with the cluster QC aggregator smart
// contract. This contract is deployed to the service account as part of a
// collection of smart contracts that facilitate and manage epoch transitions.
type QCContractClient interface {

	// SubmitVote submits the given vote to the cluster QC aggregator smart
	// contract. This function returns only once the transaction has been
	// processed by the network. An error is returned if the transaction has
	// failed and should be re-submitted.
	// Error returns:
	//   - network.TransientError for any errors from the underlying client, if the retry period has been exceeded
	//   - errTransactionExpired if the transaction has expired
	//   - errTransactionReverted if the transaction execution reverted
	//   - generic error in case of unexpected critical failure
	SubmitVote(ctx context.Context, vote *model.Vote) error

	// Voted returns true if we have successfully submitted a vote to the
	// cluster QC aggregator smart contract for the current epoch.
	// Error returns:
	//   - network.TransientError for any errors from the underlying Flow client
	//   - generic error in case of unexpected critical failures
	Voted(ctx context.Context) (bool, error)
}

// EpochLookup enables looking up epochs by view.
// Only Epochs that are fully committed can be retrieved.
// CAUTION: EpochLookup should only be used for querying the previous, current, or next epoch.
type EpochLookup interface {

	// EpochForView returns the counter of the epoch that the input view belongs to.
	// Note: The EpochLookup component processes EpochExtended notifications which will
	// extend the view range for the latest epoch.
	//
	// Returns model.ErrViewForUnknownEpoch if the input does not fall within the range of a known epoch.
	EpochForView(view uint64) (epochCounter uint64, err error)
}
