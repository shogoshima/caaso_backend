package models

type UserInput struct {
	ID          string `json:"id" binding:"required"`
	Nusp        string `json:"nusp" binding:"required"`
	DisplayName string `json:"displayName" binding:"required"`
	Email       string `json:"email" binding:"required"`
	PhotoUrl    string `json:"photoUrl" binding:"required"`
}
