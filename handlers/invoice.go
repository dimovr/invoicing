package handlers

import (
	"net/http"
	"strconv"
	"todo-item-app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InvoiceHandler struct {
	DB *gorm.DB
}

func NewInvoiceHandler(db *gorm.DB) *InvoiceHandler {
	return &InvoiceHandler{DB: db}
}

func (h *InvoiceHandler) GetInvoices(c *gin.Context) {
	var invoices []models.Invoice
	h.DB.Preload("Supplier").Preload("LineItems.Item").Find(&invoices)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"invoices": invoices,
		"active":   "invoices",
		"Title":    "Invoices",
	})
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

func (h *InvoiceHandler) GetInvoicesPartial(c *gin.Context) {
	var invoices []models.Invoice
	h.DB.Preload("Supplier").Preload("LineItems.Item").Find(&invoices)
	c.HTML(http.StatusOK, "invoice.html", gin.H{
		"invoices": invoices,
	})
}

func (h *InvoiceHandler) GetInvoiceForm(c *gin.Context) {
	// Fetch suppliers for the dropdown
	var suppliers []models.Supplier
	h.DB.Find(&suppliers)

	// Fetch items for the line item dropdowns
	var items []models.Item
	h.DB.Find(&items)

	// Get the current state of line items from query parameters
	lineItems := []map[string]interface{}{}
	lineItemCountStr := c.Query("lineItemCount")
	lineItemCount, err := strconv.Atoi(lineItemCountStr)
	if err != nil || lineItemCount <= 0 {
		lineItemCount = 1 // Default to 1 line item
	}

	for i := 0; i < lineItemCount; i++ {
		itemIDStr := c.Query("line_items[" + strconv.Itoa(i) + "][item_id]")
		countStr := c.Query("line_items[" + strconv.Itoa(i) + "][count]")
		include := c.Query("include_" + strconv.Itoa(i))
		if include == "false" {
			continue // Skip this line item if it was removed
		}

		itemID, _ := strconv.Atoi(itemIDStr)
		count, _ := strconv.Atoi(countStr)
		if count <= 0 {
			count = 1 // Default to 1 if count is invalid
		}

		// Fetch item details if an item is selected
		var item models.Item
		var itemPrice float64
		var itemUnit string
		var itemTaxRate int
		if itemID > 0 {
			if err := h.DB.First(&item, itemID).Error; err == nil {
				itemPrice = item.Price
				itemUnit = item.Unit
				itemTaxRate = item.TaxRate
			}
		} else {
			itemPrice = 0.00
			itemUnit = "-"
			itemTaxRate = 0
		}

		// Calculate totals for this line item
		subtotal := itemPrice * float64(count)
		taxAmount := subtotal * float64(itemTaxRate) / 100
		total := subtotal + taxAmount

		lineItems = append(lineItems, map[string]interface{}{
			"Index":     i,
			"ItemID":    itemID,
			"Count":     count,
			"Price":     itemPrice,
			"Unit":      itemUnit,
			"TaxRate":   itemTaxRate,
			"Subtotal":  subtotal,
			"TaxAmount": taxAmount,
			"Total":     total,
		})
	}

	// Calculate invoice totals
	var invoiceSubtotal, invoiceTaxAmount, invoiceTotal float64
	for _, lineItem := range lineItems {
		invoiceSubtotal += lineItem["Subtotal"].(float64)
		invoiceTaxAmount += lineItem["TaxAmount"].(float64)
		invoiceTotal += lineItem["Total"].(float64)
	}

	// If editing an existing invoice
	id := c.Query("id")
	if id == "" {
		c.HTML(http.StatusOK, "invoice-form.html", gin.H{
			"suppliers":        suppliers,
			"items":            items,
			"lineItems":        lineItems,
			"lineItemCount":    lineItemCount,
			"invoiceSubtotal":  invoiceSubtotal,
			"invoiceTaxAmount": invoiceTaxAmount,
			"invoiceTotal":     invoiceTotal,
		})
		return
	}

	invoiceID, err := strconv.Atoi(id)
	if err != nil || invoiceID <= 0 {
		c.String(http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	var invoice models.Invoice
	if err := h.DB.Preload("Supplier").Preload("LineItems.Item").First(&invoice, invoiceID).Error; err != nil {
		c.String(http.StatusNotFound, "Invoice not found")
		return
	}
	c.HTML(http.StatusOK, "invoice-form.html", gin.H{
		"invoice":          invoice,
		"suppliers":        suppliers,
		"items":            items,
		"lineItems":        lineItems,
		"lineItemCount":    lineItemCount,
		"invoiceSubtotal":  invoiceSubtotal,
		"invoiceTaxAmount": invoiceTaxAmount,
		"invoiceTotal":     invoiceTotal,
	})
}

func (h *InvoiceHandler) GetItems(c *gin.Context) {
	var items []models.Item
	h.DB.Find(&items)

	index := c.Query("index")
	if index == "" {
		index = "0"
	}

	c.HTML(http.StatusOK, "item-options.html", gin.H{
		"items": items,
		"Index": index,
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
