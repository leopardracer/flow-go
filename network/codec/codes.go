// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package codec

import (
	"fmt"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/libp2p/message"
	"github.com/onflow/flow-go/model/messages"
)

const (
	CodeMin uint8 = iota + 1

	// consensus
	CodeBlockProposal
	CodeBlockVote

	// protocol state sync
	CodeSyncRequest
	CodeSyncResponse
	CodeRangeRequest
	CodeBatchRequest
	CodeBlockResponse

	// cluster consensus
	CodeClusterBlockProposal
	CodeClusterBlockVote
	CodeClusterBlockResponse

	// collections, guarantees & transactions
	CodeCollectionGuarantee
	CodeTransaction
	CodeTransactionBody

	// core messages for execution & verification
	CodeExecutionReceipt
	CodeResultApproval

	// execution state synchronization
	CodeExecutionStateSyncRequest
	CodeExecutionStateDelta

	// data exchange for execution of blocks
	CodeChunkDataRequest
	CodeChunkDataResponse

	// result approvals
	CodeApprovalRequest
	CodeApprovalResponse

	// generic entity exchange engines
	CodeEntityRequest
	CodeEntityResponse

	// testing
	CodeEcho

	// DKG
	CodeDKGMessage

	CodeMax
)

// MessageCodeFromInterface returns the correct Code based on the underlying type of message v.
func MessageCodeFromInterface(v interface{}) (uint8, string, error) {
	s := what(v)
	switch v.(type) {
	// consensus
	case *messages.BlockProposal:
		return CodeBlockProposal, s, nil
	case *messages.BlockVote:
		return CodeBlockVote, s, nil

	// cluster consensus
	case *messages.ClusterBlockProposal:
		return CodeClusterBlockProposal, s, nil
	case *messages.ClusterBlockVote:
		return CodeClusterBlockVote, s, nil
	case *messages.ClusterBlockResponse:
		return CodeClusterBlockResponse, s, nil

	// protocol state sync
	case *messages.SyncRequest:
		return CodeSyncRequest, s, nil
	case *messages.SyncResponse:
		return CodeSyncResponse, s, nil
	case *messages.RangeRequest:
		return CodeRangeRequest, s, nil
	case *messages.BatchRequest:
		return CodeBatchRequest, s, nil
	case *messages.BlockResponse:
		return CodeBlockResponse, s, nil

	// collections, guarantees & transactions
	case *flow.CollectionGuarantee:
		return CodeCollectionGuarantee, s, nil
	case *flow.TransactionBody:
		return CodeTransactionBody, s, nil
	case *flow.Transaction:
		return CodeTransaction, s, nil

	// core messages for execution & verification
	case *flow.ExecutionReceipt:
		return CodeExecutionReceipt, s, nil
	case *flow.ResultApproval:
		return CodeResultApproval, s, nil

	// execution state synchronization
	case *messages.ExecutionStateSyncRequest:
		return CodeExecutionStateSyncRequest, s, nil
	case *messages.ExecutionStateDelta:
		return CodeExecutionStateDelta, s, nil

	// data exchange for execution of blocks
	case *messages.ChunkDataRequest:
		return CodeChunkDataRequest, s, nil
	case *messages.ChunkDataResponse:
		return CodeChunkDataResponse, s, nil

	// result approvals
	case *messages.ApprovalRequest:
		return CodeApprovalRequest, s, nil
	case *messages.ApprovalResponse:
		return CodeApprovalResponse, s, nil

	// generic entity exchange engines
	case *messages.EntityRequest:
		return CodeEntityRequest, s, nil
	case *messages.EntityResponse:
		return CodeEntityResponse, s, nil

	// testing
	case *message.TestMessage:
		return CodeEcho, s, nil

	// dkg
	case *messages.DKGMessage:
		return CodeDKGMessage, s, nil

	default:
		return 0, "", fmt.Errorf("invalid encode type (%T)", v)
	}
}

// InterfaceFromMessageCode returns an interface with the correct underlying go type
// of the message code represents.
// Expected error returns during normal operations:
//   - ErrUnknownMsgCode if message code does not match any of the configured message codes above.
func InterfaceFromMessageCode(code uint8) (interface{}, string, error) {
	switch code {
	// consensus
	case CodeBlockProposal:
		return &messages.BlockProposal{}, what(messages.BlockProposal{}), nil
	case CodeBlockVote:
		return &messages.BlockVote{}, what(messages.BlockVote{}), nil

	// cluster consensus
	case CodeClusterBlockProposal:
		return &messages.ClusterBlockProposal{}, what(messages.ClusterBlockProposal{}), nil
	case CodeClusterBlockVote:
		return &messages.ClusterBlockVote{}, what(messages.ClusterBlockVote{}), nil
	case CodeClusterBlockResponse:
		return &messages.ClusterBlockResponse{}, what(messages.ClusterBlockResponse{}), nil

	// protocol state sync
	case CodeSyncRequest:
		return &messages.SyncRequest{}, what(messages.SyncRequest{}), nil
	case CodeSyncResponse:
		return &messages.SyncResponse{}, what(messages.SyncResponse{}), nil
	case CodeRangeRequest:
		return &messages.RangeRequest{}, what(messages.RangeRequest{}), nil
	case CodeBatchRequest:
		return &messages.BatchRequest{}, what(messages.BatchRequest{}), nil
	case CodeBlockResponse:
		return &messages.BlockResponse{}, what(messages.BlockResponse{}), nil

	// collections, guarantees & transactions
	case CodeCollectionGuarantee:
		return &flow.CollectionGuarantee{}, what(flow.CollectionGuarantee{}), nil
	case CodeTransactionBody:
		return &flow.TransactionBody{}, what(flow.TransactionBody{}), nil
	case CodeTransaction:
		return &flow.Transaction{}, what(flow.Transaction{}), nil

	// core messages for execution & verification
	case CodeExecutionReceipt:
		return &flow.ExecutionReceipt{}, what(flow.ExecutionReceipt{}), nil
	case CodeResultApproval:
		return &flow.ResultApproval{}, what(flow.ResultApproval{}), nil

	// execution state synchronization
	case CodeExecutionStateSyncRequest:
		return &messages.ExecutionStateSyncRequest{}, what(messages.ExecutionStateSyncRequest{}), nil
	case CodeExecutionStateDelta:
		return &messages.ExecutionStateDelta{}, what(messages.ExecutionStateDelta{}), nil

	// data exchange for execution of blocks
	case CodeChunkDataRequest:
		return &messages.ChunkDataRequest{}, what(messages.ChunkDataRequest{}), nil
	case CodeChunkDataResponse:
		return &messages.ChunkDataResponse{}, what(messages.ChunkDataResponse{}), nil

	// result approvals
	case CodeApprovalRequest:
		return &messages.ApprovalRequest{}, what(messages.ApprovalRequest{}), nil
	case CodeApprovalResponse:
		return &messages.ApprovalResponse{}, what(messages.ApprovalResponse{}), nil

	// generic entity exchange engines
	case CodeEntityRequest:
		return &messages.EntityRequest{}, what(messages.EntityRequest{}), nil
	case CodeEntityResponse:
		return &messages.EntityResponse{}, what(messages.EntityResponse{}), nil

	// dkg
	case CodeDKGMessage:
		return &messages.DKGMessage{}, what(messages.DKGMessage{}), nil

	// test messages
	case CodeEcho:
		return &message.TestMessage{}, what(message.TestMessage{}), nil

	default:
		return nil, "", NewUnknownMsgCodeErr(code)
	}
}

func what(v interface{}) string {
	return fmt.Sprintf("%T", v)
}
