package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

var db *mgo.Database

// Models
type User struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Name        string        `json:"name"`
	Email       string        `json:"email"`
	CompanyName string        `json:"company_name" bson:"company_name"`
	Password    string        `json:"password"`
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

// Controllers
func Register(c *gin.Context) {
	var payload User
	_ = c.BindJSON(&payload)

	if payload.Name == "" || payload.Email == "" || payload.CompanyName == "" || payload.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "All fields are required"})
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
	}

	payload.Password = string(hashedPass)
	payload.ID = bson.NewObjectId()
	_ = db.C("users").Insert(&payload)

}

func CreateInvoice(c *gin.Context) {
	var payload = struct {
		Invoice
		Transactions []Transaction `json:"transactions"`
	}{}

	c.BindJSON(&payload)

	if payload.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "All fields are required"})
		return
	}

	var invoice Invoice
	invoice.ID = bson.NewObjectId()
	invoice.Name = payload.Name
	invoice.UserID = payload.UserID
	_ = db.C("invoices").Insert(&invoice)

	for _, v := range payload.Transactions {
		v.ID = bson.NewObjectId()
		v.InvoiceID = invoice.ID.Hex()
		_ = db.C("transactions").Insert(&v)
	}
}

func GetAllInvoiceOfAUser(c *gin.Context) {
	var invoicesData = struct {
		Status       bool `json:"status"`
		Transactions []struct {
			ID        bson.ObjectId `json:"id" bson:"_id"`
			Name      string        `json:"name" bson:"name"`
			UserID    string        `json:"user_id"`
			Paid      bool          `json:"paid"`
			Price     float32       `json:"price" bson:"price"`
			InvoiceID string        `json:"invoice_id" bson:"invoice_id"`
		}
	}{}

	invoicesData.Status = true

	var invoices []Invoice
	err := db.C("invoices").Find(bson.M{"user_id": bson.M{"$eq": c.Param("userId")}}).All(&invoices)
	if err != nil {
		log.Println(err)
	}

	for _, v := range invoices {
		var Transactions []struct {
			ID        bson.ObjectId `json:"id" bson:"_id"`
			Name      string        `json:"name" bson:"name"`
			UserID    string        `json:"user_id"`
			Paid      bool          `json:"paid"`
			Price     float32       `json:"price" bson:"price"`
			InvoiceID string        `json:"invoice_id" bson:"invoice_id"`
		}

		err = db.C("transactions").Find(bson.M{"invoice_id": bson.M{"$eq": v.ID.Hex()}}).All(&Transactions)
		if err != nil {
			log.Println(err)
		}

		for _, v2 := range Transactions {
			v2.UserID = v.UserID
			v2.Paid = v.Paid

			invoicesData.Transactions = append(invoicesData.Transactions, v2)
		}
	}

	c.JSON(200, invoicesData)
}

func GetOneInvoiceOfAUser(c *gin.Context) {
	var invoicesData = struct {
		Status       bool `json:"status"`
		Transactions []struct {
			ID        bson.ObjectId `json:"id" bson:"_id"`
			Name      string        `json:"name" bson:"name"`
			UserID    string        `json:"user_id"`
			Paid      bool          `json:"paid"`
			Price     float32       `json:"price" bson:"price"`
			InvoiceID string        `json:"invoice_id" bson:"invoice_id"`
		}
	}{}

	invoicesData.Status = true

	var invoice Invoice
	andQuery := bson.M{"$and": []bson.M{{"user_id": c.Param("userId")}, {"_id": bson.ObjectIdHex(c.Param("invoiceId"))}}}
	err := db.C("invoices").Find(andQuery).One(&invoice)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(invoice)

	err = db.C("transactions").Find(bson.M{"invoice_id": bson.M{"$eq": invoice.ID.Hex()}}).All(&invoicesData.Transactions)
	if err != nil {
		log.Println(err)
	}

	for i := range invoicesData.Transactions {
		invoicesData.Transactions[i].UserID = invoice.UserID
		invoicesData.Transactions[i].Paid = invoice.Paid
	}

	c.JSON(200, invoicesData)
}

func main() {
	app := gin.Default()

	sess, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		log.Fatal(err)
	}

	db = sess.DB("invoice")

	app.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to Invoicing App")
	})

	app.POST("/register", Register)
	app.POST("/invoices", CreateInvoice)
	app.GET("/invoices/user/:userId", GetAllInvoiceOfAUser)
	app.GET("/invoices/user/:userId/:invoiceId", GetOneInvoiceOfAUser)
	_ = app.Run(":3128")
}
