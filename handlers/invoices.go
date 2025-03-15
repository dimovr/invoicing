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

	fmt.Println(len(invoice.LineItems))

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

	// Check if item exists for invoiceid if yes return
	var existingItem models.InvoiceItem
	if err := ic.DB.Where("invoice_id = ? AND item_id = ?", invoiceIDInt, itemID).First(&existingItem).Error; err == nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Proizvod je već dodat",
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

	// Get item details
	var item models.Item
	if err := ic.DB.First(&item, itemID).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"error": "Item not found",
		})
		return
	}

	// Calculate values
	valueWithoutDiscount := price * quantity
	valueWithDiscount := valueWithoutDiscount * (1 - discount/100)
	vatAmount := valueWithDiscount * (float64(item.TaxRate) / 100)

	// Create the invoice item
	invoiceItem := models.InvoiceItem{
		ItemID:          uint(itemID),
		Name:            item.Name,
		Unit:            item.Unit,
		InvoiceID:       uint(invoiceIDInt),
		Quantity:        quantity,
		Price:           price,
		PriceDifference: discount,
		Value:           valueWithDiscount,
		ValueWithoutVat: valueWithDiscount,
		VatRate:         float64(item.TaxRate),
		VatAmount:       vatAmount,
		ValueWithVat:    valueWithDiscount + vatAmount,
		UnitPrice:       price,
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

	if count == 0 {
		// If no items left, return the "no items" row
		c.HTML(http.StatusOK, "invoice-line-item-empty.html", nil)
	} else {
		// Return empty response (the row will be removed by HTMX)
		c.Status(http.StatusOK)
	}
}

// CompleteInvoice finalizes an invoice and redirects to the view page
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

	// Calculate totals
	var subtotal float64
	var taxAmount float64

	for _, item := range invoice.LineItems {
		// Calculate item values based on discount
		valueWithoutDiscount := item.Price * item.Quantity
		valueWithDiscount := valueWithoutDiscount * (1 - item.PriceDifference/100)
		vatAmount := valueWithDiscount * (item.VatRate / 100)

		subtotal += valueWithDiscount
		taxAmount += vatAmount
	}

	// Update invoice totals
	invoice.Subtotal = subtotal
	invoice.TaxAmount = taxAmount
	invoice.Total = subtotal + taxAmount

	if err := ic.DB.Save(&invoice).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Could not update invoice: " + err.Error(),
		})
		return
	}

	fmt.Println(len(invoice.LineItems))

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

	c.HTML(http.StatusOK, "kalkulacija.html", gin.H{
		"Invoice": invoice,
		"Company": company,
		"active":  "invoices",
		"Title":   "Invoice Details",
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
