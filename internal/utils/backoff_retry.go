package utils

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

type BackoffRetry struct {
	logger         *logrus.Logger
	maxRetries     int
	initialBackoff time.Duration
}

func NewBackoffRetry(maxRetries int) *BackoffRetry {
	return &BackoffRetry{
		logger:         logrus.WithField("module", "backoff_retry").Logger,
		maxRetries:     maxRetries,
		initialBackoff: time.Second,
	}
}

// RetryWithBackoff attempts to execute the provided function `fn` up to `maxRetries` times.
// If `fn` fails, it waits for a delay that increases exponentially after each attempt.
func (b *BackoffRetry) RetryWithBackoff(fn func(string) (float64, error), arg string) (float64, error) {
	var result float64
	var err error
	backoffDuration := b.initialBackoff
	for attempt := 1; attempt <= b.maxRetries; attempt++ {
		result, err = fn(arg)
		if err == nil {
			return result, nil
		}

		backoffDuration += time.Duration(attempt)
		// Log attempt and error
		b.logger.Warnf("Attempt %d failed with error: %v. Retrying in %s...\n", attempt, err, backoffDuration)
		time.Sleep(backoffDuration)
	}

	// Return the error after exhausting all attempts
	return 0, errors.New("max retries reached: last error was " + err.Error())
}
