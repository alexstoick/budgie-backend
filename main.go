package main

import (
	"fmt"
	"github.com/alexstoick/budgie-backend/controllers"
	"github.com/alexstoick/budgie-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

	if os.Getenv("DATABASE_URL") != "" {
		dbinfo = os.Getenv("DATABASE_URL")
	}

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

	router.GET("/users", controllers.IndexUsers)
	router.POST("/users", controllers.CreateUser)

	router.POST("/auth", controllers.AuthUser)
	router.GET("/verify_token", controllers.VerifyToken)
	router.POST("/renew_token", controllers.RenewToken)

	router.POST("/users/:id/payments", controllers.CreatePayment)

	router.GET("/users/:id/payments", controllers.GetUserPayments)
	router.GET("/users/:id/payments/:payment_id", controllers.GetPaymentBeneficiaries)

	port = ":" + os.Getenv("PORT")

	router.Run(port)
}
