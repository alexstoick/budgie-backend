package controllers

import (
	"github.com/alexstoick/hello/models"
	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
	"strconv"
)

func GetPaymentBeneficiaries(r render.Render, db gorm.DB, params martini.Params) {
	var payment models.Payment

	err := db.Find(&payment, params["payment_id"]).Error

	if err != nil {
		panic(err)
	}

	user_id, _ := strconv.ParseInt(params["id"], 10, 64)
	v, err := payment.SourceID.Value()
	if v != user_id {
		r.JSON(404, map[string]string{"error": "Resource not available"})
	} else {
	}
}

func CreatePayment(r render.Render, db gorm.DB, paymentCreator models.PaymentCreator, params martini.Params) {
	paymentCreator.ParseBeneficiaryIds()
	var payment models.Payment = paymentCreator.TransformToPayment()
	payment.CreateBeneficiaries(paymentCreator.BeneficiaryIds)
	payment.AddSource(db, params["id"])
	db.Create(&payment)
	r.JSON(200, payment)
}
