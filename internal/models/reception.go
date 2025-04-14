package models

import (
	"github.com/google/uuid"
	"time"
)

type Reception struct {
	ID        uuid.UUID       `json:"id"`
	PvzID     uuid.UUID       `json:"pvzId"`
	CreatedAt time.Time       `json:"dateTime"`
	Status    ReceptionStatus `json:"status"`
}

type ReceptionStatus string

const (
	Close      ReceptionStatus = "close"
	InProgress ReceptionStatus = "in_progress"
)

type ReceptionDto struct {
	Reception Reception  `json:"reception"`
	Products  *[]Product `json:"products"`
}
