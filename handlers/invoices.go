package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"todo-item-app/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InvoiceHandler manages invoice-related requests
type InvoiceHandler struct {
	DB *gorm.DB
}

// NewInvoiceHandler creates a new invoice handler
func NewInvoiceHandler(db *gorm.DB) *InvoiceHandler {
	return &InvoiceHandler{DB: db}
}

// GetInvoices renders the list of invoices
func (h *InvoiceHandler) GetInvoices(c *gin.Context) {
	var invoices []models.Invoice
	h.DB.Preload("Supplier").Find(&invoices)

	var suppliers []models.Supplier
	h.DB.Find(&suppliers)

	var items []models.Item
	h.DB.Find(&items)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"invoices":  invoices,
		"suppliers": suppliers,
		"items":     items,
		"active":    "invoices",
		"Title":     "Invoices",
	})
}

// GetInvoicesPartial renders just the invoice list partial
func (h *InvoiceHandler) GetInvoicesPartial(c *gin.Context) {
	var invoices []models.Invoice
	h.DB.Preload("Supplier").Find(&invoices)

	c.HTML(http.StatusOK, "invoices.html", gin.H{
		"invoices": invoices,
	})
}

func (h *InvoiceHandler) GetInvoiceCreateForm(c *gin.Context) {
	var suppliers []models.Supplier
	h.DB.Find(&suppliers)

	var items []models.Item
	h.DB.Find(&items)

	fmt.Println(suppliers)
	fmt.Println(items)

	c.HTML(http.StatusOK, "invoice-create-form.html", gin.H{
		"suppliers": suppliers,
		"items":     items,
		"active":    "invoices",
		"Title":     "Invoices",
	})
}

// CreateInvoice handles the creation of a new invoice
func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	// Parse form data
	supplierID, err := strconv.ParseUint(c.PostForm("SupplierID"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid supplier ID")
		return
	}

	subtotal, _ := strconv.ParseFloat(c.PostForm("Subtotal"), 64)
	taxAmount, _ := strconv.ParseFloat(c.PostForm("TaxAmount"), 64)
	total, _ := strconv.ParseFloat(c.PostForm("Total"), 64)

	// Create invoice
	invoice := models.Invoice{
		SupplierID: uint(supplierID),
		Subtotal:   subtotal,
		TaxAmount:  taxAmount,
		Total:      total,
	}

	// Start a transaction to ensure all operations succeed or fail together
	tx := h.DB.Begin()

	if err := tx.Create(&invoice).Error; err != nil {
		tx.Rollback()
		c.String(http.StatusInternalServerError, "Failed to create invoice")
		return
	}

	// Get line items from the form
	// We're expecting a list of item IDs and counts from hidden form fields
	lineItemIDs := c.PostFormArray("LineItems[ItemID]")
	lineItemCounts := c.PostFormArray("LineItems[Count]")

	if len(lineItemIDs) != len(lineItemCounts) {
		tx.Rollback()
		c.String(http.StatusBadRequest, "Line item data mismatch")
		return
	}

	// Process each line item
	for i := 0; i < len(lineItemIDs); i++ {
		itemID, err := strconv.ParseUint(lineItemIDs[i], 10, 64)
		if err != nil {
			tx.Rollback()
			c.String(http.StatusBadRequest, "Invalid item ID")
			return
		}

		count, err := strconv.Atoi(lineItemCounts[i])
		if err != nil {
			tx.Rollback()
			c.String(http.StatusBadRequest, "Invalid count")
			return
		}

		// Get the item to calculate the amounts
		var item models.Item
		if err := tx.First(&item, itemID).Error; err != nil {
			tx.Rollback()
			c.String(http.StatusNotFound, "Item not found")
			return
		}

		subtotal := item.Price * float64(count)
		taxAmount := subtotal * float64(item.TaxRate) / 100
		total := subtotal + taxAmount

		// Create line item
		lineItem := models.InvoiceLineItem{
			InvoiceID: invoice.ID,
			ItemID:    uint(itemID),
			Count:     count,
			Subtotal:  subtotal,
			TaxAmount: taxAmount,
			Total:     total,
		}

		if err := tx.Create(&lineItem).Error; err != nil {
			tx.Rollback()
			c.String(http.StatusInternalServerError, "Failed to create line item")
			return
		}
	}

	// Commit the transaction
	tx.Commit()

	// Preload supplier for the template
	h.DB.Preload("Supplier").First(&invoice, invoice.ID)

	// Return the invoice row to be added to the table
	c.HTML(http.StatusCreated, "invoice.html", invoice)
}

// DeleteInvoice handles the deletion of an invoice
func (h *InvoiceHandler) DeleteInvoice(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// Start transaction to delete invoice and its line items
	tx := h.DB.Begin()

	// Delete line items first (due to foreign key constraints)
	if err := tx.Where("invoice_id = ?", id).Delete(&models.InvoiceLineItem{}).Error; err != nil {
		tx.Rollback()
		c.String(http.StatusInternalServerError, "Failed to delete line items")
		return
	}

	// Delete the invoice
	if err := tx.Delete(&models.Invoice{}, id).Error; err != nil {
		tx.Rollback()
		c.String(http.StatusInternalServerError, "Failed to delete invoice")
		return
	}

	tx.Commit()
	c.String(http.StatusOK, "")
}

// AddLineItem adds a line item to the invoice form (temporary, not persisted yet)
func (h *InvoiceHandler) AddLineItem(c *gin.Context) {
	itemID, err := strconv.Atoi(c.PostForm("ItemID"))
	if err != nil || itemID <= 0 {
		c.String(http.StatusBadRequest, "Invalid item ID")
		return
	}

	count, err := strconv.Atoi(c.PostForm("Count"))
	if err != nil || count <= 0 {
		c.String(http.StatusBadRequest, "Invalid count")
		return
	}

	// Fetch the item details
	var item models.Item
	if err := h.DB.First(&item, itemID).Error; err != nil {
		c.String(http.StatusNotFound, "Item not found")
		return
	}

	// Calculate amounts
	subtotal := item.Price * float64(count)
	taxAmount := subtotal * float64(item.TaxRate) / 100
	total := subtotal + taxAmount

	// Generate a temporary ID for the line item in the form
	tempId := uuid.New().String()

	// Render the line item row
	c.HTML(http.StatusOK, "invoice-line-item.html", gin.H{
		"tempId":    tempId,
		"item":      item,
		"count":     count,
		"subtotal":  subtotal,
		"taxAmount": taxAmount,
		"total":     total,
	})

	// Also update the invoice summary
	// This would typically be done with another HTMX request in practice
	// but we include it here for simplicity with a swap-oob attribute
	c.Header("HX-Trigger", `{"updateInvoiceSummary": {"tempId": "`+tempId+`"}}`)
}

// RemoveLineItem removes a line item from the invoice form
func (h *InvoiceHandler) RemoveLineItem(c *gin.Context) {
	// tempId := c.Param("tempId")

	// No need to do any database operations since the line item hasn't been persisted yet
	// Just remove it from the UI and update the totals
	c.String(http.StatusOK, "")
	c.Header("HX-Trigger", `{"updateInvoiceSummary": {"removed": "true"}}`)
}

// GetInvoiceSummary calculates and returns the updated invoice summary
func (h *InvoiceHandler) GetInvoiceSummary(c *gin.Context) {
	var subtotal, taxAmount, total float64

	// In a real implementation, we would calculate these based on line items in the form
	// For now, we're simulating this with values from the request for demonstration
	subtotal, _ = strconv.ParseFloat(c.Query("subtotal"), 64)
	taxAmount, _ = strconv.ParseFloat(c.Query("taxAmount"), 64)
	total, _ = strconv.ParseFloat(c.Query("total"), 64)

	c.HTML(http.StatusOK, "invoice-summary.html", gin.H{
		"subtotal":  subtotal,
		"taxAmount": taxAmount,
		"total":     total,
	})
}
