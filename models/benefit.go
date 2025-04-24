package models

import "time"

type Benefit struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Title       string `json:"title" gorm:"unique"`
	Description string `json:"description"`
	PhotoUrl    string `json:"photoUrl"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
