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

type InvoiceController struct {
	DB *gorm.DB
}

type LineItemInput struct {
	ItemID   uint    `json:"item_id"`
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
	Discount float64 `json:"discount"`
	ItemName string  `json:"item_name"`
	VatRate  float64 `json:"vat_rate"`
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
			subtotal += item.Price * item.Quantity * (1 - item.PriceDifference/100)
			taxAmount += item.Price * item.Quantity * (1 - item.PriceDifference/100) * (item.VatRate / 100)
		}

		invoices[i].Subtotal = subtotal
		invoices[i].TaxAmount = taxAmount
		invoices[i].Total = subtotal + taxAmount
	}

	// Get suppliers for the invoice creation form
	var suppliers []models.Supplier
	if err := ic.DB.Find(&suppliers).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Failed to load suppliers: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "invoices.html", gin.H{
		"invoices":  invoices,
		"Suppliers": suppliers,
		"TodayDate": time.Now().Format("2006-01-02"),
		"active":    "invoices",
		"Title":     "Invoices",
	})
}

// InitializeInvoice creates a new invoice with basic info and redirects to the edit page
func (ic *InvoiceController) InitializeInvoice(c *gin.Context) {
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

	// Parse date
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

	// Redirect to the invoice edit page
	ic.GetInvoiceEditPage(c, invoice.ID)
}

// GetInvoiceEditPage loads the invoice edit page
func (ic *InvoiceController) GetInvoiceEditPage(c *gin.Context, invoiceID ...uint) {
	var id uint

	// Check if we're being called with an ID parameter
	if len(invoiceID) > 0 {
		id = invoiceID[0]
	} else {
		// Otherwise, get it from the URL parameter
		idParam := c.Param("id")
		idInt, err := strconv.Atoi(idParam)
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
				"error": "Invalid invoice ID",
			})
			return
		}
		id = uint(idInt)
	}

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
			invoice.LineItems[i].Note = item.Name // Store item name for display
		}
	}

	// Load all available items
	var items []models.Item
	if err := ic.DB.Find(&items).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Could not load items: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "invoice-edit.html", gin.H{
		"Invoice": invoice,
		"Items":   items,
		"active":  "invoices",
		"Title":   "Edit Invoice",
	})
}

// AddLineItem adds an item to an invoice
func (ic *InvoiceController) AddLineItem(c *gin.Context) {
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
		Note:            item.Name,
	}

	if err := ic.DB.Create(&invoiceItem).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Could not add item to invoice: " + err.Error(),
		})
		return
	}

	// Prepare data for the line item template
	lineItemData := struct {
		ItemID     uint
		InvoiceID  uint
		ItemName   string
		Price      string
		Quantity   string
		Discount   string
		Value      string
		VatAmount  string
		TotalValue string
	}{
		ItemID:     uint(itemID),
		InvoiceID:  uint(invoiceIDInt),
		ItemName:   item.Name,
		Price:      fmt.Sprintf("%.2f", price),
		Quantity:   fmt.Sprintf("%.2f", quantity),
		Discount:   fmt.Sprintf("%.2f", discount),
		Value:      fmt.Sprintf("%.2f", valueWithDiscount),
		VatAmount:  fmt.Sprintf("%.2f", vatAmount),
		TotalValue: fmt.Sprintf("%.2f", valueWithDiscount+vatAmount),
	}

	// Remove the "no items" row if it exists (using hx-swap-oob)
	noItemsRow := `<tr id="no-items" hx-swap-oob="true"></tr>`

	// Render the line item template
	c.HTML(http.StatusOK, "invoice-line-item.html", gin.H{
		"LineItem": lineItemData,
		"NoItems":  noItemsRow,
	})
}

// RemoveLineItem removes an item from an invoice
func (ic *InvoiceController) RemoveLineItem(c *gin.Context) {
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
func (ic *InvoiceController) CompleteInvoice(c *gin.Context) {
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

		// Update the line item
		item.Value = valueWithDiscount
		item.ValueWithoutVat = valueWithDiscount
		item.VatAmount = vatAmount
		item.ValueWithVat = valueWithDiscount + vatAmount

		if err := ic.DB.Save(&item).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
				"error": "Could not update line item: " + err.Error(),
			})
			return
		}

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

	// Redirect to the view page
	c.Redirect(http.StatusFound, fmt.Sprintf("/invoices/%d/view", invoiceIDInt))
}

// GetInvoiceDetails shows the view page for an invoice
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
			invoice.LineItems[i].Note = item.Name // Store item name for display
		}
	}

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
