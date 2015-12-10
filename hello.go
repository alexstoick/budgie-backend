package main

import (
	"encoding/json"
	"fmt"
	"github.com/alexstoick/hello/controllers"
	"github.com/alexstoick/hello/models"
	"github.com/go-martini/martini"
	//	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
)

var db gorm.DB

func connectToDb() {
	const (
		DB_USER     = ""
		DB_PASSWORD = ""
		DB_NAME     = "budgie_development"
	)
	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", DB_NAME)
	var err error
	db, err = gorm.Open("postgres", dbinfo)
	handleError(err)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var payments []models.Payment

	err := db.Preload("Source").Preload("Beneficiaries").Find(&payments).Error
	handleError(err)

	jsonResult, err := json.Marshal(payments)
	handleError(err)

	fmt.Fprintf(w, "%s", jsonResult)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func autoMigrateModels() {
	db.AutoMigrate(&models.Payment{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Beneficiary{})
}

func main() {
	connectToDb()
	autoMigrateModels()

	// Declaration loop
	models.InitPaymentOnUserType()

	m := martini.Classic()

	m.Use(render.Renderer(render.Options{}))
	m.Map(db)
	m.Get("/users", controllers.IndexUsers)
	//m.Get("/users/:id", controllers.ShowUser)
	m.Post(
		"/users/:id/payments",
		binding.Form(models.PaymentCreator{}),
		controllers.CreatePayment,
	)

	m.Get("/users/:id/payments", controllers.GetUserPayments)

	m.Get("/users/:id/payments/:payment_id", controllers.GetPaymentBeneficiaries)
	m.Run()
}
