package models

type Revenue struct {
	ID     int     `json:"id" gorm:"primaryKey;type:int"`
	Amount float64 `json:"amount"`
	Month  string  `json:"month"`
	Year   int     `json:"year"`
}
