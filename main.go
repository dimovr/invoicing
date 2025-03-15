package main

import (
	"invoicing-item-app/handlers"
	"invoicing-item-app/models"
	"text/template"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("invoicing.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Company{}, &models.Item{}, &models.Supplier{}, &models.InvoiceItem{}, &models.Invoice{})

	r := gin.Default()
	// Define custom template functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}
	// Load templates with the custom function map
	r.SetFuncMap(funcMap)
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	companyHandler := handlers.NewCompanyHandler(db)
	r.GET("/company", companyHandler.GetCompany)
	r.POST("/company", companyHandler.UpsertCompany)

	itemHandler := handlers.NewItemHandler(db)
	r.GET("/items", itemHandler.GetItems)
	r.GET("/items/list", itemHandler.GetItemsPartial)
	r.GET("/items/form", itemHandler.GetItemCreateForm)
	r.POST("/items", itemHandler.CreateItem)
	r.GET("/items/:id/edit", itemHandler.GetItemEditForm)
	r.PUT("/items/:id", itemHandler.UpdateItem)
	r.DELETE("/items/:id", itemHandler.DeleteItem)

	supplierHandler := handlers.NewSupplierHandler(db)
	r.GET("/suppliers", supplierHandler.GetSuppliers)
	r.GET("/suppliers/list", supplierHandler.GetSuppliersPartial)
	r.GET("/suppliers/form", supplierHandler.GetSupplierCreateForm)
	r.POST("/suppliers", supplierHandler.CreateSupplier)
	r.GET("/suppliers/:id/edit", supplierHandler.GetSupplierEditForm)
	r.PUT("/suppliers/:id", supplierHandler.UpdateSupplier)
	r.DELETE("/suppliers/:id", supplierHandler.DeleteSupplier)

	setupRoutes(r, db)

	r.Run(":8080")
}

func setupRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize controllers
	invoiceController := handlers.NewInvoiceController(db)

	router.GET("/invoices", invoiceController.GetInvoices)
	router.POST("/invoices/initialize", invoiceController.InitializeInvoice)
	router.POST("/invoices/:id/add-item", invoiceController.AddLineItem)
	router.DELETE("/invoices/:id/items/:item_id", invoiceController.RemoveLineItem)
	router.POST("/invoices/:id/complete", invoiceController.CompleteInvoice)
	router.GET("/invoices/:id/view", invoiceController.GetInvoiceDetails)
	router.DELETE("/invoices/:id", invoiceController.DeleteInvoice)
}
