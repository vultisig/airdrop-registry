package models

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	JobDate         time.Time `gorm:"type:date;not null"`
	Multiplier      int64
	CurrentID       int64
	CurrentVaultID  uint
	IsSuccess       bool
	IsVolumeFetched bool `gorm:"type:boolean;default:false"`
}

func (*Job) TableName() string {
	return "jobs"
}

func (j *Job) Date() string {
	return j.JobDate.Format("2006-01-02")
}
func (j *Job) DaysSince() int64 {
	return int64(time.Since(j.JobDate).Hours()) / 24
}

// start and end of epoch will be used to fetch volume data
// start of epoch is the date of job.date - multiplier days
func StartOfEpoch(jobDate time.Time) int64 {
	// fetch date of job.date , remove hour,minutes and seconds  add multiplier days
	startOfEpoch := jobDate.Truncate(24 * time.Hour).Add(time.Duration(-1 * 24 * time.Hour))
	return startOfEpoch.Unix()
}

// end of epoch is the start date of job.date
func EndOfEpoch(jobDate time.Time) int64 {
	// fetch date of job.date , remove hour,minutes and seconds  add multiplier days
	endOfEpoch := jobDate.Truncate(24 * time.Hour)
	return endOfEpoch.Unix()
}
