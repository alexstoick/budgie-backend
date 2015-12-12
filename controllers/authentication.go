package controllers

import (
	"fmt"
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/alexstoick/budgie-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"os"
	"time"
)

func RenewToken(c *gin.Context) {
	token := c.PostForm("token")

	token_byte := []byte(token)
	result, err := jws.ParseJWT(token_byte)
	if err != nil {
		c.JSON(400, map[string]interface{}{"error": "cannot parse token"})
		return
	}

	issuedAt, _ := result.Claims().IssuedAt()
	fmt.Fprintf(os.Stdout, "%f\n", issuedAt)
	result.Claims().SetIssuedAt(float64(time.Now().Unix()))
	result.Claims().SetExpiration(float64(time.Now().AddDate(0, 0, 1).Unix()))

	var key = []byte(os.Getenv("JWT_SECRET"))

	valid := result.Validate(key, crypto.SigningMethodHS512) == nil
	if !valid {
		c.JSON(401, map[string]interface{}{"error": "invalid signature"})
	} else {
		issuedAt, _ := result.Claims().IssuedAt()
		claims := jws.Claims(result.Claims())
		newJWT := jws.NewJWT(claims, crypto.SigningMethodHS512)
		fmt.Fprintf(os.Stdout, "%f\n", issuedAt)
		serialized_res, _ := newJWT.Serialize(key)
		c.JSON(200, map[string]interface{}{"token": string(serialized_res)})
	}
}

func VerifyToken(c *gin.Context) {
	token := c.Query("token")

	token_byte := []byte(token)
	result, _ := jws.ParseJWT(token_byte)

	var key = []byte(os.Getenv("JWT_SECRET"))

	valid := result.Validate(key, crypto.SigningMethodHS512) == nil

	c.JSON(200, map[string]interface{}{"valid": valid})
}

func AuthUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var testuser models.User
	testuser.Username = username
	testuser.HashPassword(password)

	var user models.User

	fake_db, _ := c.Get("db")
	db := fake_db.(gorm.DB)

	db.Find(&user, models.User{Username: username})

	if user.IsMatchingPassword(password) {
		c.JSON(200, map[string]interface{}{"token": user.GenerateJWT()})
	} else {
		c.JSON(401, map[string]interface{}{"error": "Wrong password"})
	}
}
