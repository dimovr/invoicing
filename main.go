package main

import (
	"text/template"
	"todo-item-app/handlers"
	"todo-item-app/models"

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
	db.AutoMigrate(&models.Company{}, &models.Item{}, &models.Invoice{}, &models.InvoiceLineItem{}, &models.Supplier{})

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

	// Add invoice routes
	invoiceHandler := handlers.NewInvoiceHandler(db)
	r.GET("/invoices", invoiceHandler.GetInvoices)
	r.POST("/invoices", invoiceHandler.CreateInvoice)
	r.POST("/invoices/:invoiceID/line-items", invoiceHandler.AddLineItem)
	r.DELETE("/invoices/:invoiceID/line-items/:lineItemID", invoiceHandler.RemoveLineItem)
	r.GET("/invoices/:invoiceID/summary", invoiceHandler.GetSummary)
	r.POST("/invoices/:invoiceID/finalize", invoiceHandler.FinalizeInvoice)
	r.DELETE("/invoices/:invoiceID", invoiceHandler.DeleteInvoice)

	r.Run(":8080")
}
