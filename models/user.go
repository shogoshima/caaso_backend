package models

import (
	"time"
)

type User struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	Nusp           string    `json:"nusp" gorm:"unique"`
	DisplayName    string    `json:"displayName"`
	Email          string    `json:"email" gorm:"unique"`
	PhotoUrl       string    `json:"photoUrl"`
	IsSubscribed   bool      `json:"isSubscribed" gorm:"default:false"`
	Type           string    `json:"type" gorm:"default:'Nenhum'"`
	ExpirationDate time.Time `json:"expirationDate" gorm:"default:null"`
	NuspModified   bool      `json:"nuspModified" gorm:"default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
