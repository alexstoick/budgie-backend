package controllers

import (
	"github.com/alexstoick/budgie-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

func GetPaymentBeneficiaries(c *gin.Context) {
	var payment models.Payment

	fake_db, _ := c.Get("db")
	db := fake_db.(gorm.DB)

	err := db.Find(&payment, c.Param("payment_id")).Error

	if err != nil {
		panic(err)
	}

	user_id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	v, err := payment.SourceID.Value()
	if v != user_id {
		c.JSON(404, gin.H{"error": "Resource not available"})
	} else {
		c.JSON(200, payment)
	}
}

func CreatePayment(c *gin.Context) {
	var paymentCreator models.PaymentCreator
	c.BindJSON(&paymentCreator)
	payment := paymentCreator.TransformToPayment()
	payment.CreateBeneficiaries(paymentCreator.PaymentDetails)
	fake_db, _ := c.Get("db")
	db := fake_db.(gorm.DB)
	payment.AddSource(db, c.Param("id"))
	db.Create(&payment)
	c.JSON(200, payment)
}
