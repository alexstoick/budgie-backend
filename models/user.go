package models

import (
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type User struct {
	ID             int       `json:"id" gorm:"primary_key"`
	FirstName      string    `form:"firstName" json:"firstName"`
	LastName       string    `form:"lastName" json:"lastName"`
	Username       string    `form:"username" json:"username"`
	Email          string    `form:"email"  json:"email"`
	HashedPassword string    `json:"-"`
	Payments       []Payment `json:"-"`
}

type UserForm struct {
	Username string `form: "username"`
	Password string `form: "password"`
}

func (user *User) HashPassword(password string) {
	result, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user.HashedPassword = string(result)
}

func (user *User) IsMatchingPassword(password string) bool {
	byte_existing := []byte(user.HashedPassword)
	byte_new := []byte(password)
	result := bcrypt.CompareHashAndPassword(byte_existing, byte_new)

	return result == nil
}

func (user *User) GenerateJWT() string {
	claims := jws.Claims{}
	claims.Set("scope", []string{"value", "payments"})
	claims.Set("userId", user.ID)
	claims.SetIssuedAt(float64(time.Now().Unix()))
	claims.SetExpiration(float64(time.Now().AddDate(0, 0, 1).Unix()))
	claims.SetIssuer("api.staging")
	claims.SetSubject("authentication")

	token := jws.NewJWT(claims, crypto.SigningMethodHS512)

	secret := os.Getenv("JWT_SECRET")
	key := []byte(secret)
	serialized_res, _ := token.Serialize(key)
	return string(serialized_res)
}
