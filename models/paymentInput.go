package models

import (
	"encoding/json"
	"fmt"
)

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
		return "Monthly"
	case Yearly:
		return "Yearly"
	default:
		return "Unknown"
	}
}

// –– MarshalJSON turns the enum into its String() in JSON ––
func (u UserTypes) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}
func (p PlanTypes) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

// –– UnmarshalJSON parses the JSON string back into the enum ––
func (u *UserTypes) UnmarshalJSON(data []byte) error {
	var s any
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch v := s.(type) {
	case float64: // incoming number
		*u = UserTypes(int(v))
	case string: // incoming string
		switch v {
		case "Alojamento":
			*u = Aloja
		case "Graduacao":
			*u = Grad
		case "PosGraduacao":
			*u = PostGrad
		case "Outros":
			*u = Other
		default:
			return fmt.Errorf("invalid UserType string: %q", v)
		}
	default:
		return fmt.Errorf("invalid UserType type: %T", v)
	}

	return nil
}

func (p *PlanTypes) UnmarshalJSON(data []byte) error {
	var s any
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch v := s.(type) {
	case float64: // incoming number
		*p = PlanTypes(int(v))
	case string: // incoming string
		switch v {
		case "Monthly":
			*p = Monthly
		case "Yearly":
			*p = Yearly
		default:
			return fmt.Errorf("invalid PlanType string: %q", v)
		}
	default:
		return fmt.Errorf("invalid PlanType type: %T", v)
	}

	return nil
}

type PaymentInput struct {
	UserType UserTypes `json:"userType"`
	PlanType PlanTypes `json:"planType"`
}
