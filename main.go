package main

import (
	"fmt"
	"github.com/alexstoick/budgie-backend/controllers"
	"github.com/alexstoick/budgie-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"os"
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

	fmt.Fprintf(os.Stdout, "fhjdsjfhds")
	if os.Getenv("DATABASE_URL") != "" {
		fmt.Fprintf(os.Stdout, "fhjdsjfhds")
		url := os.Getenv("DATABASE_URL")
		dbinfo, _ := pq.ParseURL(url)
		dbinfo += " sslmode=disable"
	}
	fmt.Fprintf(os.Stdout, "%s", dbinfo)

	db, err = gorm.Open("postgres", dbinfo)
	handleError(err)
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

func DatabaseMapper(c *gin.Context) {
	c.Set("db", db)

	c.Next()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	connectToDb()
	autoMigrateModels()
	godotenv.Load()

	router := gin.Default()

	router.Use(DatabaseMapper)
	router.Use(CORSMiddleware())

	v1 := router.Group("v1/")
	{
		v1.GET("/users", controllers.IndexUsers)
		v1.POST("/users", controllers.CreateUser)

		v1.POST("/login", controllers.AuthUser)
		v1.GET("/verify_token", controllers.VerifyToken)
		v1.POST("/renew_token", controllers.RenewToken)

		authentication := v1.Use(controllers.ValidateAuthentication)

		authentication.POST("/users/me/payments", controllers.CreatePayment)

		authentication.GET("/users/me/payments", controllers.GetUserPayments)
		authentication.GET("/users/me/payments/:payment_id", controllers.GetPaymentBeneficiaries)
	}

	port := ":" + os.Getenv("PORT")

	router.Run(port)
}
