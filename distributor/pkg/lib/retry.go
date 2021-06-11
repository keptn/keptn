package lib

import (
	"context"
	"fmt"
	"time"
)

//TODO: move to go-utils

type Option func(*RetryConfiguration)

func NumberOfRetries(n uint) Option {
	return func(c *RetryConfiguration) {
		c.numberOfRetries = n
	}
}

func DelayBetweenRetries(d time.Duration) Option {
	return func(c *RetryConfiguration) {
		c.delayBetweenRetries = d
	}
}

type RetryConfiguration struct {
	numberOfRetries     uint
	delayBetweenRetries time.Duration
}

type RetryFunc func() error

// Retry executes the retryFunc repeatedly until it was successful or canceled by the context
// The default number of retries is 20 and the default delay between retries is 2 seconds
func Retry(context context.Context, retryFunc RetryFunc, opts ...Option) error {
	configuration := &RetryConfiguration{numberOfRetries: 20, delayBetweenRetries: time.Second * 2}
	for _, opt := range opts {
		opt(configuration)
	}

	var i uint
	for i < configuration.numberOfRetries {
		err := retryFunc()
		if err != nil {
			select {
			case <-time.After(configuration.delayBetweenRetries):
			case <-context.Done():
				return fmt.Errorf("retry cancelled")
			}
		} else {
			return nil
		}
		i++
	}
	return fmt.Errorf("operation unsuccessful after %d retry", i)
}
