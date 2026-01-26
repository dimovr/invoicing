package main

import (
	"flag"
	"html/template"
	"invoicing-item-app/csv"
	"invoicing-item-app/handlers"
	"invoicing-item-app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	populate := flag.Bool("populate", false, "Populate database from CSV file")
	flag.Parse()

	db := initDb()

	if *populate {
		csv.Populate("artikli.csv")
	}

	r := gin.Default()

	setupTemplates(r)

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
	r.GET("/items/export", itemHandler.ExportItems)

	supplierHandler := handlers.NewSupplierHandler(db)
	r.GET("/suppliers", supplierHandler.GetSuppliers)
	r.GET("/suppliers/list", supplierHandler.GetSuppliersPartial)
	r.GET("/suppliers/form", supplierHandler.GetSupplierCreateForm)
	r.POST("/suppliers", supplierHandler.CreateSupplier)
	r.GET("/suppliers/:id/edit", supplierHandler.GetSupplierEditForm)
	r.PUT("/suppliers/:id", supplierHandler.UpdateSupplier)
	r.DELETE("/suppliers/:id", supplierHandler.DeleteSupplier)

	invoiceHandler := handlers.NewInvoiceHandler(db)
	r.GET("/invoices", invoiceHandler.GetInvoices)
	r.POST("/invoices", invoiceHandler.InitializeInvoice)
	r.POST("/invoices/:id/items", invoiceHandler.AddLineItem)
	r.DELETE("/invoices/:id/items/:item_id", invoiceHandler.RemoveLineItem)
	r.POST("/invoices/:id/complete", invoiceHandler.CompleteInvoice)
	r.GET("/invoices/:id/view", invoiceHandler.GetInvoiceDetails)
	r.GET("/invoices/:id/edit", invoiceHandler.GetInvoiceEditPage)
	r.DELETE("/invoices/:id", invoiceHandler.DeleteInvoice)

	_ = r.Run(":8080")
}

func initDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("invoicing.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	_ = db.AutoMigrate(&models.Company{}, &models.Item{}, &models.Supplier{}, &models.InvoiceItem{}, &models.Invoice{})
	return db
}

func setupTemplates(r *gin.Engine) {
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b float64) float64 {
			return a - b
		},
	}

	r.SetFuncMap(funcMap)
	r.LoadHTMLGlob("templates/*")
}
