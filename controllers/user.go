package controllers

import (
	"github.com/alexstoick/hello/models"
	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
)

func IndexUsers(r render.Render, db gorm.DB) {
	var users []models.User
	err := db.Find(&users).Error

	if err != nil {
		panic(err)
	}

	r.JSON(200, users)
}

func GetUserPayments(r render.Render, db gorm.DB, params martini.Params) {
	var payments []models.Payment
	var user models.User

	err := db.Find(&user, params["id"]).Error

	err = db.Model(&user).Preload("Beneficiaries").Preload("Beneficiaries.Beneficiary").Related(&payments, "SourceID").Error

	if err != nil {
		panic(err)
	}

	r.JSON(200, payments)
}
