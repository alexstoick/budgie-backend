package controllers

import (
	"github.com/alexstoick/budgie-backend/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/alexstoick/budgie-backend/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/alexstoick/budgie-backend/models"
)

func IndexUsers(c *gin.Context) {
	var users []models.User
	fake_db, _ := c.Get("db")
	db := fake_db.(gorm.DB)
	err := db.Find(&users).Error

	if err != nil {
		panic(err)
	}

	c.JSON(200, users)
}

func GetUserPayments(c *gin.Context) {
	var payments []models.Payment
	var user models.User

	fake_db, _ := c.Get("db")
	db := fake_db.(gorm.DB)

	err := db.Find(&user, c.Param("id")).Error

	err = db.Model(&user).Preload("Beneficiaries").Preload("Beneficiaries.Beneficiary").Related(&payments, "SourceID").Error

	if err != nil {
		panic(err)
	}

	c.JSON(200, payments)
}
