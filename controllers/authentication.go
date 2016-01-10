package controllers

import (
	"fmt"
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/alexstoick/budgie-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"os"
	"strings"
	"time"
)

type JSONResponse struct {
	statuscode int
	body       gin.H
}

var TokenParseError JSONResponse = JSONResponse{400, gin.H{"error": "Cannot Parse Token"}}
var TokenSignatureError JSONResponse = JSONResponse{401, gin.H{"error": "Invalid Signature"}}
var WrongPasswordError JSONResponse = JSONResponse{401, gin.H{"error": "Wrong Password"}}
var WrongTokenError JSONResponse = JSONResponse{403, gin.H{"error": "Wrong Token"}}

func RenewToken(c *gin.Context) {
	token := c.PostForm("token")

	token_byte := []byte(token)
	result, err := jws.ParseJWT(token_byte)
	if err != nil {
		c.JSON(TokenParseError.statuscode, TokenParseError.body)
		c.Abort()
		return
	}

	issuedAt, _ := result.Claims().IssuedAt()
	fmt.Fprintf(os.Stdout, "%f\n", issuedAt)
	result.Claims().SetIssuedAt(float64(time.Now().Unix()))
	result.Claims().SetExpiration(float64(time.Now().AddDate(0, 0, 1).Unix()))

	var key = []byte(os.Getenv("JWT_SECRET"))

	valid := result.Validate(key, crypto.SigningMethodHS512) == nil
	if !valid {
		c.JSON(TokenSignatureError.statuscode, TokenSignatureError.body)
		c.Abort()
		return
	}

	issuedAt, _ = result.Claims().IssuedAt()
	claims := jws.Claims(result.Claims())
	newJWT := jws.NewJWT(claims, crypto.SigningMethodHS512)
	fmt.Fprintf(os.Stdout, "%f\n", issuedAt)
	serialized_res, _ := newJWT.Serialize(key)
	c.JSON(200, gin.H{"token": string(serialized_res)})
}

func VerifyToken(c *gin.Context) {
	token := c.Query("token")

	token_byte := []byte(token)
	result, _ := jws.ParseJWT(token_byte)

	var key = []byte(os.Getenv("JWT_SECRET"))

	valid := result.Validate(key, crypto.SigningMethodHS512) == nil

	c.JSON(200, gin.H{"valid": valid})
}

func AuthUser(c *gin.Context) {

	var userForm models.UserForm

	c.BindJSON(&userForm)

	var user models.User

	fake_db, _ := c.Get("db")
	db := fake_db.(gorm.DB)

	db.Find(&user, models.User{Username: userForm.Username})

	if !user.IsMatchingPassword(userForm.Password) {
		c.JSON(WrongPasswordError.statuscode, WrongPasswordError.body)
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"token": user.GenerateJWT()})
}

func ValidateAuthentication(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	split_header := strings.Split(header, " ")

	var token string
	fmt.Fprintf(os.Stdout, "\nsplit_header --- %v\n", split_header)
	if len(split_header) == 2 {
		token = split_header[1]
	} else {
		token = c.Query("token")
	}
	fmt.Fprintf(os.Stdout, "\ntoken --- %v\n", token)

	token_byte := []byte(token)
	result, err := jws.ParseJWT(token_byte)

	if err != nil {
		c.JSON(TokenParseError.statuscode, TokenParseError.body)
		c.Abort()
		return
	}
	var key = []byte(os.Getenv("JWT_SECRET"))

	valid := result.Validate(key, crypto.SigningMethodHS512) == nil

	if !valid {
		c.JSON(TokenSignatureError.statuscode, TokenSignatureError.body)
		c.Abort()
		return
	}
	c.Set("userId", result.Claims().Get("userId"))

	c.Next()
}
