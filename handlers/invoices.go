package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"invoicing-item-app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InvoiceHandler struct {
	DB *gorm.DB
}

func NewInvoiceHandler(db *gorm.DB) *InvoiceHandler {
	return &InvoiceHandler{DB: db}
}

func (ic *InvoiceHandler) GetInvoices(c *gin.Context) {
	var invoices []models.Invoice
	if err := ic.DB.Preload("Supplier").Preload("LineItems").Find(&invoices).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Failed to load invoices: " + err.Error(),
		})
		return
	}

	// Get suppliers for the invoice creation form
	var suppliers []models.Supplier
	if err := ic.DB.Find(&suppliers).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Failed to load suppliers: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"invoices":  invoices,
		"suppliers": suppliers,
		"TodayDate": time.Now().Format("2006-01-02"),
		"active":    "invoices",
		"Title":     "Invoices",
	})
}

func (ic *InvoiceHandler) InitializeInvoice(c *gin.Context) {
	supplierID, err := strconv.Atoi(c.PostForm("supplier_id"))
	if err != nil || supplierID == 0 {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Please select a supplier",
		})
		return
	}

	documentNumber := c.PostForm("document_number")
	if documentNumber == "" {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Please enter a document number",
		})
		return
	}

	dateStr := c.PostForm("date")
	var invoiceDate time.Time
	if dateStr != "" {
		var err error
		fmt.Println("dateStr:", dateStr)
		invoiceDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			invoiceDate = time.Now()
		}
	} else {
		invoiceDate = time.Now()
	}

	// Create a new invoice with basic info
	invoice := models.Invoice{
		SupplierID:     uint(supplierID),
		DocumentNumber: documentNumber,
		Date:           invoiceDate,
		Subtotal:       0,
		TaxAmount:      0,
		Total:          0,
	}

	if err := ic.DB.Create(&invoice).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Could not create invoice: " + err.Error(),
		})
		return
	}

	// Redirect to the invoice edit page instead of rendering directly
	c.Redirect(http.StatusFound, fmt.Sprintf("/invoices/%d/edit", invoice.ID))
}

// GetInvoiceEditPage loads the invoice edit page
func (ic *InvoiceHandler) GetInvoiceEditPage(c *gin.Context) {
	invoiceID := c.Param("id")
	id, err := strconv.Atoi(invoiceID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Invalid invoice ID",
		})
		return
	}

	var invoice models.Invoice
	if err := ic.DB.Preload("Supplier").Preload("LineItems").First(&invoice, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"error": "Invoice not found",
		})
		return
	}

	// Load all available items
	var items []models.Item
	if err := ic.DB.Find(&items).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Could not load items: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "invoice-form.html", gin.H{
		"Invoice": invoice,
		"Items":   items,
		"active":  "invoices",
		"Title":   "Edit Invoice",
	})
}

func (ic *InvoiceHandler) AddLineItem(c *gin.Context) {
	invoiceID := c.Param("id")
	invoiceIDInt, err := strconv.Atoi(invoiceID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Invalid invoice ID",
		})
		return
	}

	// Parse form data
	itemID, err := strconv.Atoi(c.PostForm("item_id"))
	if err != nil || itemID == 0 {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Please select an item",
		})
		return
	}

	var existingItem models.InvoiceItem
	if err := ic.DB.Where("invoice_id = ? AND item_id = ?", invoiceIDInt, itemID).First(&existingItem).Error; err == nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Item already added to invoice",
		})
		return
	}

	price, err := strconv.ParseFloat(c.PostForm("price"), 64)
	if err != nil || price <= 0 {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Please enter a valid price",
		})
		return
	}

	quantity, err := strconv.ParseFloat(c.PostForm("quantity"), 64)
	if err != nil || quantity <= 0 {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Please enter a valid quantity",
		})
		return
	}

	discount, err := strconv.ParseFloat(c.PostForm("discount"), 64)
	if err != nil {
		discount = 0
	}

	var item models.Item
	if err := ic.DB.First(&item, itemID).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"error": "Item not found",
		})
		return
	}

	// Calculate values
	buyingPrice := price * (1 - discount/100)
	buyingPriceWithTax := buyingPrice * (1 + float64(item.TaxRate)/100)
	buyingSubtotal := buyingPriceWithTax * quantity
	sellingPrice := item.Price
	sellingPriceWithTax := sellingPrice
	sellingTotal := sellingPriceWithTax * quantity

	pc := float64(100*item.TaxRate) / float64(100+item.TaxRate) / 100
	taxAmount := sellingTotal * pc

	// Create the invoice item
	invoiceItem := models.InvoiceItem{
		InvoiceID:    uint(invoiceIDInt),
		ItemID:       uint(itemID),
		Name:         item.Name,
		Unit:         item.Unit,
		TaxRate:      float64(item.TaxRate),
		Discount:     discount,
		Quantity:     quantity,
		BuyingPrice:  buyingPrice,
		Subtotal:     buyingSubtotal,
		TaxAmount:    taxAmount,
		SellingPrice: sellingPrice,
		Total:        sellingTotal,
	}

	if err := ic.DB.Create(&invoiceItem).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Could not add item to invoice: " + err.Error(),
		})
		return
	}

	// Render the line item template
	c.HTML(http.StatusOK, "invoice-line-item.html", invoiceItem)
}

// RemoveLineItem removes an item from an invoice
func (ic *InvoiceHandler) RemoveLineItem(c *gin.Context) {
	invoiceID := c.Param("id")
	itemID := c.Param("item_id")

	if err := ic.DB.Where("invoice_id = ? AND item_id = ?", invoiceID, itemID).Delete(&models.InvoiceItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not remove item",
		})
		return
	}

	// Check if this was the last item
	var count int64
	ic.DB.Model(&models.InvoiceItem{}).Where("invoice_id = ?", invoiceID).Count(&count)

	c.Status(http.StatusOK)
}

func (ic *InvoiceHandler) CompleteInvoice(c *gin.Context) {
	invoiceID := c.Param("id")
	invoiceIDInt, err := strconv.Atoi(invoiceID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Invalid invoice ID",
		})
		return
	}

	var invoice models.Invoice
	if err := ic.DB.Preload("LineItems").First(&invoice, invoiceID).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"error": "Invoice not found",
		})
		return
	}

	if len(invoice.LineItems) == 0 {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Invoice has no items",
		})
		return
	}

	invoice.Subtotal = 0
	invoice.TaxAmount = 0
	invoice.Total = 0

	for _, item := range invoice.LineItems {
		invoice.Subtotal += item.Subtotal
		invoice.TaxAmount += item.TaxAmount
		invoice.Total += item.Total
	}

	if err := ic.DB.Save(&invoice).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Could not update invoice: " + err.Error(),
		})
		return
	}

	// Redirect to the view page
	c.Redirect(http.StatusFound, fmt.Sprintf("/invoices/%d/view", invoiceIDInt))
}

// GetInvoiceDetails shows the view page for an invoice
func (ic *InvoiceHandler) GetInvoiceDetails(c *gin.Context) {
	id := c.Param("id")

	var company models.Company
	if err := ic.DB.First(&company).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Could not load company: " + err.Error(),
		})
		return
	}

	var invoice models.Invoice
	if err := ic.DB.Preload("Supplier").Preload("LineItems").First(&invoice, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"error": "Invoice not found",
		})
		return
	}

	c.HTML(http.StatusOK, "invoice-full.html", gin.H{
		"Invoice":   invoice,
		"Company":   company,
		"TodayDate": time.Now().Format("02.01.2006"),
		"active":    "invoices",
		"Title":     "Invoice Details",
	})
}

func (ic *InvoiceHandler) DeleteInvoice(c *gin.Context) {
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

	c.String(http.StatusOK, "")
}
