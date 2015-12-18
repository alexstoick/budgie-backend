package models

import "database/sql"

type Beneficiary struct {
	ID            int           `json:"-" gorm:"primary_key"`
	Payment       Payment       `json:"-"`
	PaymentID     sql.NullInt64 `json:"-"`
	Beneficiary   User          `json:"user,omitempty"`
	BeneficiaryID sql.NullInt64 `json:"-"`
	Amount        float64       `json:"amount"`
}
