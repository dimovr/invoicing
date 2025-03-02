package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"todo-item-app/models"

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
	fmt.Println("invoices")
	var invoices []models.Invoice
	ic.DB.Find(&invoices)
	fmt.Println(invoices)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"invoices": invoices,
		"active":   "invoices",
		"Title":    "Invoices",
	})
}

// GetInvoiceForm renders the invoice creation form
func (ic *InvoiceController) GetInvoiceForm(c *gin.Context) {
	// Get current company
	var company models.Company
	if err := ic.DB.First(&company).Error; err != nil {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": "Company not found",
		})
		return
	}

	// Get all items
	var items []models.Item
	if err := ic.DB.Find(&items).Error; err != nil {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": "Could not load items",
		})
		return
	}

	// Get all suppliers
	var suppliers []models.Supplier
	if err := ic.DB.Find(&suppliers).Error; err != nil {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": "Could not load suppliers",
		})
		return
	}

	// Format current date
	currentDate := time.Now().Format("02.01.2006")

	c.HTML(http.StatusOK, "kalkulacija.html", gin.H{
		"Company":     company,
		"Items":       items,
		"Suppliers":   suppliers,
		"CurrentDate": currentDate,
	})
}

// GetItemDetails fetches details for a specific item
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

	// HTMX response is handled by JavaScript in the template
	c.Status(http.StatusOK)
}

// SaveInvoice saves a new invoice
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

	// documentType := c.PostForm("document_type")
	documentNumber := c.PostForm("document_number")
	if documentNumber == "" {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": "Unesite broj dokumenta",
		})
		return
	}

	// Parse items JSON
	var items []models.InvoiceItem
	itemsJSON := c.PostForm("items_json")
	if err := json.Unmarshal([]byte(itemsJSON), &items); err != nil {
		fmt.Println("Greska tokom obrade artikala:", err)
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

	// Here you would save the invoice to your database
	// For example:
	invoice := models.Invoice{
		SupplierID:     uint(supplierID),
		DocumentNumber: documentNumber,
		Date:           time.Now(),
		// Add other fields as needed
	}

	if err := ic.DB.Create(&invoice).Error; err != nil {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": "Could not save invoice: " + err.Error(),
		})
		return
	}

	// Save invoice items
	for _, item := range items {
		invoiceItem := models.InvoiceItem{
			InvoiceID:       invoice.ID,
			ItemID:          item.ItemID,
			Quantity:        item.Quantity,
			Price:           item.Price,
			DependentCosts:  item.DependentCosts,
			PriceDifference: item.PriceDifference,
			VatRate:         item.VatRate,
			Note:            item.Note,
			// Add other fields as needed
		}

		if err := ic.DB.Create(&invoiceItem).Error; err != nil {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"error": "Could not save invoice item: " + err.Error(),
			})
			return
		}
	}

	// Success message
	c.String(http.StatusOK, `<div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded" role="alert">
		<strong>Success!</strong> Invoice saved successfully.
	</div>`)
}
