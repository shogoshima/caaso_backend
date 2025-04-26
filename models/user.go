package models

import (
	"time"
)

type User struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	Token          string    `json:"token" gorm:"type:uuid;default:gen_random_uuid();unique"`
	DisplayName    string    `json:"displayName"`
	Email          string    `json:"email" gorm:"unique"`
	PhotoUrl       string    `json:"photoUrl"`
	IsSubscribed   bool      `json:"isSubscribed" gorm:"default:false"`
	Type           string    `json:"type" gorm:"default:'Nenhum'"`
	ExpirationDate time.Time `json:"expirationDate" gorm:"default:null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
