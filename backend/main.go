package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"log"
)

var db *mgo.Database


func main() {
	app := gin.Default()

	app.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           86400,
	}))

	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		log.Fatal(err)
	}

	db = session.DB("invoice")

	app.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to Invoicing App")
	})

	app.POST("/register", Register)
	app.POST("/login", Login)
	app.POST("/invoices", CreateInvoice)
	app.GET("/invoices/user/:userId", GetAllInvoiceOfAUser)
	app.GET("/invoices/user/:userId/:invoiceId", GetOneInvoiceOfAUser)
	_ = app.Run(":3128")
}
