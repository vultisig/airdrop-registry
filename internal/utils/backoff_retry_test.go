package utils

import (
	"errors"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// TestRetryWithBackoff_Success tests the function when it succeeds within the retry limit.
func TestRetryWithBackoff_Success(t *testing.T) {
	// Arrange
	logger := logrus.New()
	retries := 3
	val := ""
	backoff := NewBackoffRetry(retries)
	backoff.logger = logger

	// Mock function to succeed on the second attempt
	attemptCount := 0
	mockFn := func(arg string) (float64, error) {
		attemptCount++
		if attemptCount == 2 {
			val = arg
			return 42.0, nil // Success on the second try
		}
		return 0, errors.New("temporary error")
	}

	// Act
	result, err := backoff.RetryWithBackoff(mockFn, "param")

	// Assert
	assert.NoError(t, err, "Expected no error after successful retry")
	assert.Equal(t, 42.0, result, "Expected result to match the successful return value")
	assert.Equal(t, 2, attemptCount, "Expected function to succeed on the second attempt")
	assert.Equal(t, "param", val, "Expected function to receive the correct argument")
}

// TestRetryWithBackoff_Failure tests the function when it fails after exhausting all retries.
func TestRetryWithBackoff_Failure(t *testing.T) {
	// Arrange
	logger := logrus.New()
	retries := 2
	backoff := NewBackoffRetry(retries)
	backoff.logger = logger

	// Mock function to always fail
	mockFn := func(arg string) (float64, error) {
		return 0, errors.New("persistent error")
	}

	// Act
	startTime := time.Now()
	result, err := backoff.RetryWithBackoff(mockFn, "")
	elapsedTime := time.Since(startTime)

	// Assert
	assert.Error(t, err, "Expected an error after exhausting retries")
	assert.Equal(t, "max retries reached: last error was persistent error", err.Error())
	assert.Equal(t, 0.0, result, "Expected result to be zero after failure")
	assert.GreaterOrEqual(t, elapsedTime, time.Duration(retries)*backoff.initialBackoff, "Expected total retry time to exceed cumulative backoff delay")
}
