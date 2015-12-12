package models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

type Payment struct {
	ID            int           `json:"id" gorm:"primary_key"`
	Amount        float64       `json:"amount"`
	Source        User          `json:"-"`
	SourceID      sql.NullInt64 `json:"-"`
	Beneficiaries []Beneficiary `json:"beneficiaries,omitempty"`
}

type PaymentCreator struct {
	Amount               float64 `form:"amount" json:"amount"`
	BeneficiaryIdsString string  `form:"beneficiary_ids" json:"dsadsadsa"`
	BeneficiaryIds       []int64 `form:"-" json:"beneficiary_ids"`
	SourceId             string  `form:"-"`
}

func (payment *PaymentCreator) ParseBeneficiaryIds() {
	splits := strings.Split(payment.BeneficiaryIdsString, ",")

	for _, split := range splits {
		converted, _ := strconv.ParseInt(split, 10, 64)
		payment.BeneficiaryIds = append(payment.BeneficiaryIds, converted)
	}
}

func (payment *Payment) AddSource(db gorm.DB, source_id string) {
	db.Find(&payment.Source, source_id)
}

func (paymentCreator PaymentCreator) TransformToPayment() Payment {
	p := Payment{
		Amount: paymentCreator.Amount,
	}
	return p
}

func (payment *Payment) CreateBeneficiaries(BeneficiaryIds []int64) {
	count := len(BeneficiaryIds)
	for _, beneficiary_id := range BeneficiaryIds {
		payment.Beneficiaries = append(
			payment.Beneficiaries,
			Beneficiary{
				BeneficiaryID: sql.NullInt64{Int64: beneficiary_id, Valid: true},
				PaymentID:     sql.NullInt64{Int64: int64(payment.ID), Valid: true},
				Amount:        payment.Amount / float64(count),
			},
		)
	}
}
