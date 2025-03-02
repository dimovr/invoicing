package handlers

import (
	"net/http"
	"strconv"
	"todo-item-app/models"

	"github.com/gin-gonic/gin"
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

	c.HTML(http.StatusOK, "index.html", gin.H{
		"invoices":  invoices,
		"suppliers": suppliers,
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

// CreateInvoice - POST /invoices
func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	supplierIDStr := c.PostForm("supplierID")
	supplierID, err := strconv.Atoi(supplierIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid supplierID")
		return
	}

	invoice := models.Invoice{SupplierID: uint(supplierID), Status: "draft"}
	if err := h.DB.Create(&invoice).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to create invoice")
		return
	}

	// Redirect to the invoice management page
	c.Redirect(http.StatusFound, "/invoices/"+strconv.Itoa(int(invoice.ID)))
}

// ShowInvoice - GET /invoices/:invoiceID
func (h *InvoiceHandler) ShowInvoice(c *gin.Context) {
	invoiceIDStr := c.Param("invoiceID")
	invoiceID, err := strconv.Atoi(invoiceIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid invoiceID")
		return
	}

	var invoice models.Invoice
	if err := h.DB.First(&invoice, invoiceID).Error; err != nil {
		c.String(http.StatusNotFound, "Invoice not found")
		return
	}

	var items []models.Item
	h.DB.Find(&items)

	c.HTML(http.StatusOK, "invoice_content.html", gin.H{
		"InvoiceID": invoice.ID,
		"Items":     items,
	})
}

func (h *InvoiceHandler) AddLineItem(c *gin.Context) {
	invoiceIDStr := c.Param("invoiceID")
	invoiceID, err := strconv.Atoi(invoiceIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid invoiceID")
		return
	}

	var invoice models.Invoice
	if err := h.DB.First(&invoice, invoiceID).Error; err != nil || invoice.Status != "draft" {
		c.String(http.StatusNotFound, "Invoice not found or not draft")
		return
	}

	itemIDStr := c.PostForm("itemId")
	quantityStr := c.PostForm("quantity")
	sellingPriceStr := c.PostForm("sellingPrice")

	itemID, _ := strconv.Atoi(itemIDStr)
	quantity, _ := strconv.Atoi(quantityStr)
	sellingPrice, _ := strconv.ParseFloat(sellingPriceStr, 64)

	var item models.Item
	if err := h.DB.First(&item, itemID).Error; err != nil {
		c.String(http.StatusBadRequest, "Invalid itemID")
		return
	}

	subtotal := float64(quantity) * sellingPrice
	lineItem := models.LineItem{
		InvoiceID:    uint(invoiceID),
		ItemID:       uint(itemID),
		Quantity:     quantity,
		SellingPrice: sellingPrice,
		Subtotal:     subtotal,
	}
	if err := h.DB.Create(&lineItem).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to add line item")
		return
	}

	c.HTML(http.StatusOK, "line_item_row.html", gin.H{
		"LineItemID":   lineItem.ID,
		"InvoiceID":    invoiceID,
		"Name":         item.Name,
		"Unit":         item.Unit,
		"Quantity":     quantity,
		"SellingPrice": sellingPrice,
		"Subtotal":     subtotal,
	})
}

func (h *InvoiceHandler) RemoveLineItem(c *gin.Context) {
	invoiceIDStr := c.Param("invoiceID")
	lineItemIDStr := c.Param("lineItemID")

	invoiceID, _ := strconv.Atoi(invoiceIDStr)
	lineItemID, _ := strconv.Atoi(lineItemIDStr)

	var invoice models.Invoice
	if err := h.DB.First(&invoice, invoiceID).Error; err != nil || invoice.Status != "draft" {
		c.String(http.StatusNotFound, "Invoice not found or not draft")
		return
	}

	if err := h.DB.Delete(&models.LineItem{}, lineItemID).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to remove line item")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *InvoiceHandler) GetSummary(c *gin.Context) {
	invoiceIDStr := c.Param("invoiceID")
	invoiceID, err := strconv.Atoi(invoiceIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid invoiceID")
		return
	}

	var lineItems []models.LineItem
	h.DB.Where("invoice_id = ?", invoiceID).Find(&lineItems)

	subtotal := 0.0
	for _, li := range lineItems {
		subtotal += li.Subtotal
	}
	total := subtotal

	c.HTML(http.StatusOK, "invoice_summary.html", gin.H{
		"InvoiceID": invoiceID,
		"Subtotal":  subtotal,
		"Total":     total,
	})
}

func (h *InvoiceHandler) FinalizeInvoice(c *gin.Context) {
	invoiceIDStr := c.Param("invoiceID")
	invoiceID, err := strconv.Atoi(invoiceIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid invoiceID")
		return
	}

	var invoice models.Invoice
	if err := h.DB.First(&invoice, invoiceID).Error; err != nil || invoice.Status != "draft" {
		c.String(http.StatusNotFound, "Invoice not found or not draft")
		return
	}

	if err := h.DB.Model(&invoice).Update("status", "finalized").Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to finalize invoice")
		return
	}

	c.HTML(http.StatusOK, "invoice_finalized.html", gin.H{
		"InvoiceID": invoiceID,
	})
}

func (h *InvoiceHandler) DeleteInvoice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("invoiceID"))
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	var invoice models.Invoice
	if err := h.DB.First(&invoice, id).Error; err != nil {
		c.String(http.StatusNotFound, "Invoice not found")
		return
	}
	h.DB.Delete(&invoice)
	c.String(http.StatusOK, "")
}
