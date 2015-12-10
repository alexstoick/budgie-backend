package models

import (
	"github.com/graphql-go/graphql"
)

type User struct {
	ID       int       `json:"id" gorm:"primary_key"`
	Name     string    `json:"name"`
	Payments []Payment `json:"-"`
}

var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
	},
})

func InitPaymentOnUserType() {
	UserType.AddFieldConfig("payments", &graphql.Field{
		Type: graphql.NewList(paymentType),
	})
}
