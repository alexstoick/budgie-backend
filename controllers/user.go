package controllers

import (
	"github.com/alexstoick/budgie-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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

func CreateUser(c *gin.Context) {
	var user models.User

	c.Bind(&user)
	user.HashPassword(c.PostForm("password"))

	fake_db, _ := c.Get("db")
	db := fake_db.(gorm.DB)
	db.Create(&user)
	c.JSON(200, user)
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
