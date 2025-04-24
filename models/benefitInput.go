package models

type BenefitInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	PhotoUrl    string `json:"photoUrl" binding:"required"`
}
