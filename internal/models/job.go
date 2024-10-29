package models

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	JobDate        time.Time `gorm:"type:date;not null"`
	Multiplier     int64
	CurrentID      int64
	CurrentVaultID uint
	IsSuccess      bool
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
