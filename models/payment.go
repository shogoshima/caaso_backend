package models

import (
	"time"
)

type Payment struct {
	ID           int       `json:"id" gorm:"primaryKey;type:int"`
	UserId       string    `json:"userId" gorm:"index"`
	QrCode       string    `json:"qrCode"`
	Amount       float64   `json:"amount"`
	IsPaid       bool      `json:"isPaid" gorm:"default:false"`
	DateApproved time.Time `json:"dateApproved" gorm:"default:null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
