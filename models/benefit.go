package models

type Benefit struct {
	ID          uint   `json:"id,omitempty" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"unique"`
	Description string `json:"description"`
	PhotoUrl    string `json:"photoUrl"`
}
