package client

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"

	"google.golang.org/grpc"

	"github.com/dapperlabs/bamboo-emulator/crypto"
	"github.com/dapperlabs/bamboo-emulator/data"
	"github.com/dapperlabs/bamboo-emulator/gen/grpc/services/accessv1"
)

type Client struct {
	conn       *grpc.ClientConn
	grpcClient accessv1.BambooAccessAPIClient
}

func New(host string, port int) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	grpcClient := accessv1.NewBambooAccessAPIClient(conn)

	return &Client{
		conn:       conn,
		grpcClient: grpcClient,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

// GetBlockByHash fetches a block by hash.
func (c *Client) GetBlockByHash(ctx context.Context, h crypto.Hash) (*data.Block, error) {
	res, err := c.grpcClient.GetBlockByHash(
		ctx,
		&accessv1.GetBlockByHashRequest{Hash: h.Bytes()},
	)
	if err != nil {
		return nil, err
	}

	block := res.GetBlock()

	timestamp, err := ptypes.Timestamp(block.GetTimestamp())
	if err != nil {
		return nil, err
	}

	return &data.Block{
		Number:            block.GetNumber(),
		PrevBlockHash:     crypto.BytesToHash(block.GetPrevBlockHash()),
		Timestamp:         timestamp,
		Status:            data.BlockStatus(block.GetStatus()),
		TransactionHashes: crypto.BytesToHashes(block.GetTransactionHashes()),
	}, nil
}

// GetBlockByNumber fetches a block by number.
func (c *Client) GetBlockByNumber(ctx context.Context, n uint64) (*data.Block, error) {
	res, err := c.grpcClient.GetBlockByNumber(
		ctx,
		&accessv1.GetBlockByNumberRequest{Number: n},
	)
	if err != nil {
		return nil, err
	}

	block := res.GetBlock()

	timestamp, err := ptypes.Timestamp(block.GetTimestamp())
	if err != nil {
		return nil, err
	}

	return &data.Block{
		Number:            block.GetNumber(),
		PrevBlockHash:     crypto.BytesToHash(block.GetPrevBlockHash()),
		Timestamp:         timestamp,
		Status:            data.BlockStatus(block.GetStatus()),
		TransactionHashes: crypto.BytesToHashes(block.GetTransactionHashes()),
	}, nil
}

// LogCommands displays all the usable commands to a client.
func LogCommands() {
	// TODO: log all help commands available to Client
	fmt.Println("here are all the commands you can use!")
}

// CreateAccount creates an account for the client.
func CreateAccount() {
	// TODO: create a new account in the client's wallet
}
