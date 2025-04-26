package models

import "time"

type UserUpdateInput struct {
	DisplayName    string    `json:"displayName" binding:"required"`
	Email          string    `json:"email" binding:"required"`
	IsSubscribed   bool      `json:"isSubscribed"`
	ExpirationDate time.Time `json:"expirationDate" binding:"required"`
}
