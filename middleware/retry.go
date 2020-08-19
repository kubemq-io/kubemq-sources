package middleware

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/pkg/retry"
	"github.com/kubemq-hub/kubemq-sources/types"
	"math"
	"time"
)

var delayTypeMap = map[string]string{
	"back-off": "back-off",
	"fixed":    "fixed",
	"random":   "random",
	"":         "",
}

type RetryMiddleware struct {
	opts []retry.Option
}

func parseRetryOptions(meta types.Metadata) ([]retry.Option, error) {
	var opts []retry.Option
	attempts, err := meta.ParseIntWithRange("retry_attempts", 1, 1, math.MaxInt32)
	if err != nil {
		return nil, fmt.Errorf("invalid retry attempts value")
	}
	opts = append(opts, retry.Attempts(uint(attempts)))

	delayMilliseconds, err := meta.ParseIntWithRange("retry_delay_milliseconds", 100, 0, math.MaxInt32)
	if err != nil {
		return nil, fmt.Errorf("invalid retry delay millisecond svalue")
	}
	opts = append(opts, retry.MaxDelay(time.Duration(delayMilliseconds)*time.Millisecond))

	maxJitterMilliseconds, err := meta.ParseIntWithRange("retry_max_jitter_milliseconds", 100, 1, math.MaxInt32)
	if err != nil {
		return nil, fmt.Errorf("invalid retry delay jitter millisecond value")
	}
	opts = append(opts, retry.MaxJitter(time.Duration(maxJitterMilliseconds)*time.Millisecond))

	delayType, err := meta.ParseStringMap("retry_delay_type", delayTypeMap)
	if err != nil {
		return nil, fmt.Errorf("invalid retry delay type value")
	}
	switch delayType {
	case "back-off", "":
		opts = append(opts, retry.DelayType(retry.BackOffDelay))
	case "fixed":
		opts = append(opts, retry.DelayType(retry.FixedDelay))
	case "random":
		opts = append(opts, retry.DelayType(retry.RandomDelay))
	}
	return opts, nil
}

func NewRetryMiddleware(meta types.Metadata, log *logger.Logger) (*RetryMiddleware, error) {
	opts, err := parseRetryOptions(meta)
	if err != nil {
		return nil, fmt.Errorf("error parsing retry options, %w", err)
	}
	if log != nil {
		opts = append(opts, retry.OnRetry(func(n uint, err error) {
			log.Errorf("retry %d failed, error: %s", n, err.Error())
		}))
	}
	return &RetryMiddleware{
		opts: opts,
	}, nil
}
