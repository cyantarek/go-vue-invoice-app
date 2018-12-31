package main

import (
	"gopkg.in/mgo.v2/bson"
)

// Models
type User struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Name        string        `json:"name"`
	Email       string        `json:"email"`
	CompanyName string        `json:"company_name" bson:"company_name"`
	Password    string        `json:"-"`
}

type Invoice struct {
	ID     bson.ObjectId `json:"id" bson:"_id"`
	Name   string        `json:"name"`
	Paid   bool          `json:"paid"`
	UserID string        `json:"user_id" bson:"user_id"`
}

type Transaction struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Name      string        `json:"name"`
	Price     float32       `json:"price"`
	InvoiceID string        `json:"invoice_id" bson:"invoice_id"`
}
