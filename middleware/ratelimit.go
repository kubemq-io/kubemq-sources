package middleware

import (
	"fmt"
	"math"

	"github.com/kubemq-io/kubemq-sources/pkg/ratelimit"
	"github.com/kubemq-io/kubemq-sources/types"
)

type RateLimitMiddleware struct {
	rateLimiter ratelimit.Limiter
}

func NewRateLimitMiddleware(meta types.Metadata) (*RateLimitMiddleware, error) {
	rpc, err := meta.ParseIntWithRange("rate_per_seconds", 0, 0, math.MaxInt32)
	if err != nil {
		return nil, fmt.Errorf("invalid rate limiter rate per second value, %w", err)
	}
	rl := &RateLimitMiddleware{}
	if rpc > 0 {
		rl.rateLimiter = ratelimit.New(rpc, ratelimit.WithoutSlack)
	} else {
		rl.rateLimiter = ratelimit.NewUnlimited()
	}
	return rl, nil
}

func (rl *RateLimitMiddleware) Take() {
	_ = rl.rateLimiter.Take()
}
