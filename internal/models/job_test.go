package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// add new test for Job struct
func TestJobEpoch(t *testing.T) {
	job := &Job{
		JobDate:    time.Date(2023, 10, 5, 5, 4, 3, 2, time.UTC),
		Multiplier: 1,
	}
	expectedStartOfEpoch := time.Date(2023, 10, 4, 0, 0, 0, 0, time.UTC).Unix()
	assert.Equal(t, job.StartOfEpoch(), expectedStartOfEpoch)

	expectedEndOfEpoch := time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC).Unix()
	assert.Equal(t, job.EndOfEpoch(), expectedEndOfEpoch)

	job.Multiplier = 2

	expectedStartOfEpoch = time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC).Unix()
	assert.Equal(t, job.StartOfEpoch(), expectedStartOfEpoch)
	assert.Equal(t, job.EndOfEpoch(), expectedEndOfEpoch)
}
