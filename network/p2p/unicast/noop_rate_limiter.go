package unicast

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/onflow/flow-go/network/message"
)

type NoopRateLimiter struct{}

func (n *NoopRateLimiter) Allow(_ peer.ID, _ *message.Message) bool {
	return true
}

func (n *NoopRateLimiter) IsRateLimited(_ peer.ID) bool {
	return false
}

func (n *NoopRateLimiter) SetTimeNowFunc(_ GetTimeNow) {}

func (n *NoopRateLimiter) Stop() {}

func (n *NoopRateLimiter) Start() {}

// NoopRateLimiters returns noop rate limiters.
func NoopRateLimiters() *RateLimiters {
	return &RateLimiters{
		MessageRateLimiter:   &NoopRateLimiter{},
		BandWidthRateLimiter: &NoopRateLimiter{},
		OnRateLimitedPeer:    nil,
		dryRun:               true,
	}
}
