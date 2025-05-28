package models

type AlojaUser struct {
	ID    uint   `json:"id,omitempty" gorm:"primaryKey"`
	Email string `json:"email" gorm:"unique;not null"`
}
