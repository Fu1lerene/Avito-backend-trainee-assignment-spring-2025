package models

import "time"

type Filter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	Limit     int
}
