package models

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID          uuid.UUID   `json:"id"`
	ReceptionID uuid.UUID   `json:"receptionId"`
	Type        ProductType `json:"type"`
	CreatedAt   time.Time   `json:"dateTime"`
}

type ProductType string

const (
	Shoes       ProductType = "обувь"
	Clothes     ProductType = "одежда"
	Electronics ProductType = "электроника"
)
