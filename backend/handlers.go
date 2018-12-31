package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

// Controllers
func Register(c *gin.Context) {
	var newUser User
	_ = c.Bind(&newUser)

	if newUser.Name == "" || newUser.Email == "" || newUser.CompanyName == "" || newUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "All fields are required"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
	}

	newUser.Password = string(hashedPassword)
	newUser.ID = bson.NewObjectId()
	_ = db.C("users").Insert(&newUser)
	c.JSON(200, gin.H{"token": uuid.New(), "user": newUser, "status": true})

}

func Login(c *gin.Context) {
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if payload.Email == "" && payload.Password == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "All fields are required"})
		return
	}
	var user User
	err := db.C("users").Find(bson.M{"email": payload.Email}).One(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password not match"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": uuid.New(), "user": user})
}

func CreateInvoice(c *gin.Context) {
	payload := struct {
		Name         string `json:"name"`
		Transactions []struct {
			ID        bson.ObjectId `bson:"_id"`
			Name      string        `json:"name"`
			InvoiceID string        `json:"invoice_id" bson:"invoice_id"`
			Price     int           `json:"price"`
		} `json:"transactions"`
		UserID string `json:"user_id"`
	}{}
	_ = c.BindJSON(&payload)

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
	var invoices []Invoice
	err := db.C("invoices").Find(bson.M{"user_id": bson.M{"$eq": c.Param("userId")}}).All(&invoices)
	if err != nil {
		log.Println(err)
	}

	c.JSON(200, gin.H{"status": true, "user_id": c.Param("userId"), "invoices": invoices})
}

func GetOneInvoiceOfAUser(c *gin.Context) {
	var invoice Invoice
	andQuery := bson.M{"$and": []bson.M{{"user_id": c.Param("userId")}, {"_id": bson.ObjectIdHex(c.Param("invoiceId"))}}}
	err := db.C("invoices").Find(andQuery).One(&invoice)
	if err != nil {
		log.Println(err)
	}

	var transactions []Transaction
	err = db.C("transactions").Find(bson.M{"invoice_id": bson.M{"$eq": invoice.ID.Hex()}}).All(&transactions)
	if err != nil {
		log.Println(err)
	}

	c.JSON(200, gin.H{"status": true, "invoice_id": invoice.ID, "user_id": c.Param("userId"), "transactions": transactions, "paid": invoice.Paid})
}
