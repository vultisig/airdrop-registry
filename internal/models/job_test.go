package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// add new test for Job struct
func TestJobEpoch(t *testing.T) {
	job := &Job{
		JobDate: time.Date(2023, 10, 5, 5, 4, 3, 2, time.UTC),
	}
	expectedStartOfEpoch := time.Date(2023, 10, 4, 0, 0, 0, 0, time.UTC).Unix()
	assert.Equal(t, StartOfEpoch(job.JobDate), expectedStartOfEpoch)

	expectedEndOfEpoch := time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC).Unix()
	assert.Equal(t, EndOfEpoch(job.JobDate), expectedEndOfEpoch)
}
