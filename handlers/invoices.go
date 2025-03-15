package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"invoicing-item-app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InvoiceController struct {
	DB *gorm.DB
}

func NewInvoiceController(db *gorm.DB) *InvoiceController {
	return &InvoiceController{DB: db}
}

func (ic *InvoiceController) GetInvoices(c *gin.Context) {
	var invoices []models.Invoice
	if err := ic.DB.Preload("Supplier").Preload("LineItems").Find(&invoices).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Failed to load invoices: " + err.Error(),
		})
		return
	}

	// Calculate totals for each invoice
	for i := range invoices {
		var subtotal float64
		var taxAmount float64

		for _, item := range invoices[i].LineItems {
			subtotal += item.Price * item.Quantity
			taxAmount += item.Price * item.Quantity * (item.VatRate / 100)
		}

		invoices[i].Subtotal = subtotal
		invoices[i].TaxAmount = taxAmount
		invoices[i].Total = subtotal + taxAmount
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"invoices": invoices,
		"active":   "invoices",
		"Title":    "Invoices",
	})
}

func (ic *InvoiceController) GetInvoiceForm(c *gin.Context) {
	var company models.Company
	if err := ic.DB.First(&company).Error; err != nil {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": "Company not found",
		})
		return
	}

	var items []models.Item
	if err := ic.DB.Find(&items).Error; err != nil {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": "Could not load items",
		})
		return
	}

	var suppliers []models.Supplier
	if err := ic.DB.Find(&suppliers).Error; err != nil {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": "Could not load suppliers",
		})
		return
	}

	currentDate := time.Now().Format("02.01.2006")

	c.HTML(http.StatusOK, "invoice-create-form.html", gin.H{
		"Company":     company,
		"Items":       items,
		"Suppliers":   suppliers,
		"CurrentDate": currentDate,
		"active":      "create-invoice",
		"Title":       "Create Invoice",
	})
}

func (ic *InvoiceController) GetItemDetails(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Query("itemId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var item models.Item
	if err := ic.DB.First(&item, itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Return item details as JSON
	c.JSON(http.StatusOK, item)
}

func (ic *InvoiceController) SaveInvoice(c *gin.Context) {
	_, err := strconv.Atoi(c.PostForm("company_id"))
	if err != nil {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": "Invalid company ID",
		})
		return
	}

	supplierID, err := strconv.Atoi(c.PostForm("supplier_id"))
	if err != nil || supplierID == 0 {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": "Izaberite dobavljača",
		})
		return
	}

	documentNumber := c.PostForm("document_number")
	if documentNumber == "" {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": "Unesite broj dokumenta",
		})
		return
	}

	// Parse date from form
	dateStr := c.PostForm("invoice_date")
	var invoiceDate time.Time
	if dateStr != "" {
		var err error
		invoiceDate, err = time.Parse("02.01.2006", dateStr)
		if err != nil {
			invoiceDate = time.Now()
		}
	} else {
		invoiceDate = time.Now()
	}

	var items []models.InvoiceItem
	itemsJSON := c.PostForm("items")
	if err := json.Unmarshal([]byte(itemsJSON), &items); err != nil {
		fmt.Println("Error processing items:", err)
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": "Greska tokom obrade artikala",
		})
		return
	}

	if len(items) == 0 {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": "Morate izabrati bar jedan artikal",
		})
		return
	}

	// Calculate totals
	var subtotal float64
	var taxAmount float64
	for _, item := range items {
		// Calculate item values
		item.ValueWithoutVat = item.Price * item.Quantity
		item.VatAmount = item.ValueWithoutVat * (item.VatRate / 100)
		item.ValueWithVat = item.ValueWithoutVat + item.VatAmount

		// Calculate unit price with dependent costs
		if item.Quantity > 0 {
			item.UnitPrice = (item.Price + (item.DependentCosts / item.Quantity)) * (1 + (item.PriceDifference / 100))
		}

		subtotal += item.ValueWithoutVat
		taxAmount += item.VatAmount
	}

	invoice := models.Invoice{
		SupplierID:     uint(supplierID),
		DocumentNumber: documentNumber,
		Date:           invoiceDate,
		Subtotal:       subtotal,
		TaxAmount:      taxAmount,
		Total:          subtotal + taxAmount,
	}

	if err := ic.DB.Create(&invoice).Error; err != nil {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": "Could not save invoice: " + err.Error(),
		})
		return
	}

	for _, item := range items {
		invoiceItem := models.InvoiceItem{
			InvoiceID:       invoice.ID,
			ItemID:          item.ItemID,
			Quantity:        item.Quantity,
			Price:           item.Price,
			DependentCosts:  item.DependentCosts,
			PriceDifference: item.PriceDifference,
			Value:           item.ValueWithoutVat,
			ValueWithoutVat: item.ValueWithoutVat,
			VatRate:         item.VatRate,
			VatAmount:       item.VatAmount,
			ValueWithVat:    item.ValueWithVat,
			UnitPrice:       item.UnitPrice,
			Note:            item.Note,
		}

		if err := ic.DB.Create(&invoiceItem).Error; err != nil {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"error": "Could not save invoice item: " + err.Error(),
			})
			return
		}
	}

	c.String(http.StatusOK, `<div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded" role="alert">
		<strong>Uspesno sacuvano</strong>
	</div>`)
}

func (ic *InvoiceController) GetInvoiceDetails(c *gin.Context) {
	id := c.Param("id")

	var invoice models.Invoice
	if err := ic.DB.Preload("Supplier").Preload("LineItems").First(&invoice, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"error": "Invoice not found",
		})
		return
	}

	// Get item details for each line item
	for i, lineItem := range invoice.LineItems {
		var item models.Item
		if err := ic.DB.First(&item, lineItem.ItemID).Error; err == nil {
			invoice.LineItems[i].Note = item.Name // Store item name in note field for display
		}
	}

	// Calculate totals
	var subtotal float64
	var taxAmount float64
	for _, item := range invoice.LineItems {
		subtotal += item.ValueWithoutVat
		taxAmount += item.VatAmount
	}

	invoice.Subtotal = subtotal
	invoice.TaxAmount = taxAmount
	invoice.Total = subtotal + taxAmount

	c.HTML(http.StatusOK, "invoice-details.html", gin.H{
		"Invoice": invoice,
		"active":  "invoices",
		"Title":   "Invoice Details",
	})
}

func (ic *InvoiceController) DeleteInvoice(c *gin.Context) {
	id := c.Param("id")

	// First delete all line items
	if err := ic.DB.Where("invoice_id = ?", id).Delete(&models.InvoiceItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete invoice items"})
		return
	}

	// Then delete the invoice
	if err := ic.DB.Delete(&models.Invoice{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete invoice"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
