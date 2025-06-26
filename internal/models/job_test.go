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
	expectedDate := time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC).Unix()
	assert.Equal(t, GetDate(job.JobDate), expectedDate)

	job = &Job{
		JobDate: time.Date(2023, 10, 5, 1, 10, 0, 0, time.FixedZone("UTC+2", 2*60*60)),
	}
	// job date in UTC is still 2023-10-4
	expectedDate = time.Date(2023, 10, 4, 0, 0, 0, 0, time.UTC).Unix()
	assert.Equal(t, GetDate(job.JobDate), expectedDate)

}
