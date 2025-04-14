package models

import (
	"github.com/google/uuid"
	"time"
)

type Pvz struct {
	ID        uuid.UUID `json:"id"`
	City      PvzCity   `json:"city"`
	CreatedAt time.Time `json:"registrationDate"`
}

type PvzCity string

const (
	Kazan           PvzCity = "Казань"
	Moscow          PvzCity = "Москва"
	SaintPetersburg PvzCity = "Санкт-Петербург"
)

type PzvDto struct {
	PVZ        Pvz             `json:"pvz"`
	Receptions *[]ReceptionDto `json:"receptions"`
}
