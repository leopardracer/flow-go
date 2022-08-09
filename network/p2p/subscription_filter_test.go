package p2p_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/id"
	"github.com/onflow/flow-go/module/irrecoverable"
	"github.com/onflow/flow-go/network/channels"
	"github.com/onflow/flow-go/network/p2p"
	"github.com/onflow/flow-go/utils/unittest"
)

// TestFilterSubscribe tests that if node X is filtered out on a specific channel by node Y's subscription
// filter, then node Y will never propagate any of node X's messages on that channel
func TestFilterSubscribe(t *testing.T) {
	// TODO: skip for now due to bug in libp2p gossipsub implementation:
	// https://github.com/libp2p/go-libp2p-pubsub/issues/449
	unittest.SkipUnless(t, unittest.TEST_TODO, "skip for now due to bug in libp2p gossipsub implementation: https://github.com/libp2p/go-libp2p-pubsub/issues/449")

	sporkId := unittest.IdentifierFixture()
	identity1, privateKey1 := unittest.IdentityWithNetworkingKeyFixture(unittest.WithRole(flow.RoleAccess))
	identity2, privateKey2 := unittest.IdentityWithNetworkingKeyFixture(unittest.WithRole(flow.RoleAccess))
	ids := flow.IdentityList{identity1, identity2}

	node1 := createNode(t, identity1.NodeID, privateKey1, sporkId, zerolog.Nop(), withSubscriptionFilter(subscriptionFilter(identity1, ids)))
	node2 := createNode(t, identity2.NodeID, privateKey2, sporkId, zerolog.Nop(), withSubscriptionFilter(subscriptionFilter(identity2, ids)))

	unstakedKey := unittest.NetworkingPrivKeyFixture()
	unstakedNode := createNode(t, flow.ZeroID, unstakedKey, sporkId, zerolog.Nop())

	require.NoError(t, node1.AddPeer(context.TODO(), *host.InfoFromHost(node2.Host())))
	require.NoError(t, node1.AddPeer(context.TODO(), *host.InfoFromHost(unstakedNode.Host())))

	badTopic := channels.TopicFromChannel(channels.SyncCommittee, sporkId)

	sub1, err := node1.Subscribe(badTopic, unittest.NetworkCodec(), unittest.AllowAllPeerFilter())
	require.NoError(t, err)

	sub2, err := node2.Subscribe(badTopic, unittest.NetworkCodec(), unittest.AllowAllPeerFilter())
	require.NoError(t, err)

	unstakedSub, err := unstakedNode.Subscribe(badTopic, unittest.NetworkCodec(), unittest.AllowAllPeerFilter())
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		return len(node1.ListPeers(badTopic.String())) > 0 &&
			len(node2.ListPeers(badTopic.String())) > 0 &&
			len(unstakedNode.ListPeers(badTopic.String())) > 0
	}, 1*time.Second, 100*time.Millisecond)

	// check that node1 and node2 don't accept unstakedNode as a peer
	require.Never(t, func() bool {
		for _, pid := range node1.ListPeers(badTopic.String()) {
			if pid == unstakedNode.Host().ID() {
				return true
			}
		}
		return false
	}, 1*time.Second, 100*time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(2)

	testPublish := func(wg *sync.WaitGroup, from *p2p.Node, sub *pubsub.Subscription) {
		data := []byte("hello")

		err := from.Publish(context.TODO(), badTopic, data)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		msg, err := sub.Next(ctx)
		cancel()
		require.NoError(t, err)
		require.Equal(t, msg.Data, data)

		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		_, err = unstakedSub.Next(ctx)
		cancel()
		require.ErrorIs(t, err, context.DeadlineExceeded)

		wg.Done()
	}

	// publish a message from node 1 and check that only node2 receives
	testPublish(&wg, node1, sub2)

	// publish a message from node 2 and check that only node1 receives
	testPublish(&wg, node2, sub1)

	unittest.RequireReturnsBefore(t, wg.Wait, 1*time.Second, "timeout performing publish test")
}

// TestCanSubscribe tests that the subscription filter blocks a node from subscribing
// to channel that its role shouldn't subscribe to
func TestCanSubscribe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	signalCtx, errChan := irrecoverable.WithSignaler(ctx)
	go unittest.NoIrrecoverableError(ctx, t, errChan)

	identity, privateKey := unittest.IdentityWithNetworkingKeyFixture(unittest.WithRole(flow.RoleCollection))
	sporkId := unittest.IdentifierFixture()

	collectionNode := createNode(t, identity.NodeID, privateKey, sporkId, zerolog.Nop(), withSubscriptionFilter(subscriptionFilter(identity, flow.IdentityList{identity})))

	startNode(t, signalCtx, collectionNode, 100*time.Millisecond)
	defer stopNode(t, collectionNode, cancel, 1*time.Second)

	goodTopic := channels.TopicFromChannel(channels.ProvideCollections, sporkId)
	_, err := collectionNode.Subscribe(goodTopic, unittest.NetworkCodec(), unittest.AllowAllPeerFilter())
	require.NoError(t, err)

	var badTopic channels.Topic
	allowedChannels := make(map[channels.Channel]struct{})
	for _, ch := range channels.ChannelsByRole(flow.RoleCollection) {
		allowedChannels[ch] = struct{}{}
	}
	for _, ch := range channels.Channels() {
		if _, ok := allowedChannels[ch]; !ok {
			badTopic = channels.TopicFromChannel(ch, sporkId)
			break
		}
	}
	_, err = collectionNode.Subscribe(badTopic, unittest.NetworkCodec(), unittest.AllowAllPeerFilter())
	require.Error(t, err)

	clusterTopic := channels.TopicFromChannel(channels.SyncCluster(flow.Emulator), sporkId)
	_, err = collectionNode.Subscribe(clusterTopic, unittest.NetworkCodec(), unittest.AllowAllPeerFilter())
	require.NoError(t, err)
}

func subscriptionFilter(self *flow.Identity, ids flow.IdentityList) pubsub.SubscriptionFilter {
	idProvider := id.NewFixedIdentityProvider(ids)
	return p2p.NewRoleBasedFilter(self.Role, idProvider)
}
