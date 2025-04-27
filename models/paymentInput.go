package models

// –– your existing types ––
type UserTypes int

const (
	Aloja UserTypes = iota
	Grad
	PostGrad
	Other
)

type PlanTypes int

const (
	Monthly PlanTypes = iota
	Yearly
)

// –– String() methods for readability ––
func (u UserTypes) String() string {
	switch u {
	case Aloja:
		return "Alojamento"
	case Grad:
		return "Graduação"
	case PostGrad:
		return "Pós-Graduação"
	case Other:
		return "Outros"
	default:
		return "Unknown"
	}
}

func (p PlanTypes) String() string {
	switch p {
	case Monthly:
		return "Mensal"
	case Yearly:
		return "Anual"
	default:
		return "Unknown"
	}
}

type PaymentInput struct {
	UserType string `json:"userType"`
	PlanType string `json:"planType"`
}
