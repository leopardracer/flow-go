package cmd

import (
	"context"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	client "github.com/onflow/flow-go-sdk/access/grpc"

	"github.com/onflow/flow-go/engine/common/rpc/convert"
	"github.com/onflow/flow-go/model/bootstrap"
	"github.com/onflow/flow-go/model/flow"
	ioutils "github.com/onflow/flow-go/utils/io"
)

// snapshotCmd represents a command to download the latest protocol state snapshot
// from an access node and write to disk
var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Download the latest protocol state snapshot from an access node and write to disk",
	Long:  `Download the latest protocol state snapshot from an access node and write to disk`,
	Run:   snapshot,
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
	addSnapshotCmdFlags()
}

func addSnapshotCmdFlags() {
	snapshotCmd.Flags().StringVarP(&flagAccessAddress, "access-address", "a", "", "the address of an access node")
	_ = snapshotCmd.MarkFlagRequired("access-address")
}

// snapshot downloads a protocol snapshot from an access node and writes it to disk
func snapshot(cmd *cobra.Command, args []string) {
	log.Info().Msg("running download snapshot")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	nodeIDString, err := readNodeID()
	if err != nil {
		log.Fatal().Err(err).Msg("could not read node ID")
	}

	nodeID, err := flow.HexStringToIdentifier(nodeIDString)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse nodeID to flow.Identifier")
	}

	// create a flow client with given access address
	flowClient, err := client.NewClient(
		flagAccessAddress,
		client.WithGRPCDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create flow client")
	}

	// get latest snapshot bytes encoded as JSON
	bytes, err := flowClient.GetLatestProtocolStateSnapshot(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get latest protocol snapshot from access node")
	}

	// unmarshal bytes to snapshot
	snapshot, err := convert.BytesToInmemSnapshot(bytes)
	if err != nil {
		log.Fatal().Err(err).Msg("could not convert array of bytes to snapshot")
	}

	// check if given NodeID is part of the current or next epoch
	currentEpoch, err := snapshot.Epochs().Current()
	if err != nil {
		log.Fatal().Err(err).Msg("could not get current epoch")
	}
	currentIdentities := currentEpoch.InitialIdentities()
	if _, exists := currentIdentities.ByNodeID(nodeID); exists {
		err := ioutils.WriteFile(filepath.Join(flagBootDir, bootstrap.PathRootProtocolStateSnapshot), bytes)
		if err != nil {
			log.Fatal().Err(err).Msg("could not write snapshot to disk")
		}
		return
	}

	nextEpoch, err := snapshot.Epochs().NextCommitted()
	if err != nil {
		log.Fatal().Err(err).Msg("could not get next committed epoch")
	}
	nextIdentities := nextEpoch.InitialIdentities()
	if _, exists := nextIdentities.ByNodeID(nodeID); exists {
		err := ioutils.WriteFile(filepath.Join(flagBootDir, bootstrap.PathRootProtocolStateSnapshot), bytes)
		if err != nil {
			log.Fatal().Err(err).Msg("could not write snapshot to disk")
		}
		return
	}

	log.Fatal().Str("node_id", nodeID.String()).Msgf("could not write snapshot, given node ID does not belong to current or next epoch")
}
