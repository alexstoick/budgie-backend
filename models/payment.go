package models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type Payment struct {
	ID            int           `json:"id" gorm:"primary_key"`
	Amount        float64       `json:"amount"`
	Source        User          `json:"-"`
	SourceID      sql.NullInt64 `json:"-"`
	Beneficiaries []Beneficiary `json:"beneficiaries,omitempty"`
}

type PaymentCreator struct {
	Amount         float64          `form:"amount" json:"amount"`
	SourceId       string           `form:"-"`
	PaymentDetails []PaymentDetails `form:"details" json:"details"`
}

type PaymentDetails struct {
	UserId int64   `form:"userId" json:"userId"`
	Amount float64 `form:"amount" json:"amount"`
}

func (payment *Payment) AddSource(db gorm.DB, source_id float64) {
	db.Find(&payment.Source, source_id)
}

func (paymentCreator PaymentCreator) TransformToPayment() Payment {
	p := Payment{
		Amount: paymentCreator.Amount,
	}
	return p
}

func (payment *Payment) CreateBeneficiaries(paymentDetails []PaymentDetails) {
	for _, paymentDetail := range paymentDetails {
		payment.Beneficiaries = append(
			payment.Beneficiaries,
			Beneficiary{
				BeneficiaryID: sql.NullInt64{Int64: paymentDetail.UserId, Valid: true},
				PaymentID:     sql.NullInt64{Int64: int64(payment.ID), Valid: true},
				Amount:        paymentDetail.Amount,
			},
		)
	}
}
