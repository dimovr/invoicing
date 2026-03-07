package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"invoicing-item-app/csv"
	"invoicing-item-app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ItemHandler struct {
	DB *gorm.DB
}

func NewItemHandler(db *gorm.DB) *ItemHandler {
	return &ItemHandler{
		DB: db,
	}
}

func (h *ItemHandler) GetItems(c *gin.Context) {
	var items []models.Item
	query := h.DB

	query.Find(&items)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"items":  items,
		"active": "items",
		"Title":  "Items",
	})
}

func (h *ItemHandler) GetItemsPartial(c *gin.Context) {
	var items []models.Item
	query := h.DB

	query.Find(&items)
	c.HTML(http.StatusOK, "items_list.html", gin.H{
		"items": items,
	})
}

func (h *ItemHandler) CreateItem(c *gin.Context) {
	var item models.Item
	if err := c.Bind(&item); err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	h.DB.Create(&item)
	c.HTML(http.StatusCreated, "item.html", item)
}

func (h *ItemHandler) DeleteItem(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var item models.Item
	if err := h.DB.First(&item, id).Error; err != nil {
		c.String(http.StatusNotFound, "Not found")
		return
	}
	h.DB.Delete(&item)
	c.String(http.StatusOK, "")
}

func (h *ItemHandler) GetItemCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "item-create-form.html", gin.H{})
}

func (h *ItemHandler) GetItemEditForm(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var item models.Item
	if err := h.DB.First(&item, id).Error; err != nil {
		c.String(http.StatusNotFound, "Not found")
		return
	}
	c.HTML(http.StatusOK, "item-edit-form.html", gin.H{"item": item})
}

func (h *ItemHandler) UpdateItem(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var item models.Item
	if err := h.DB.First(&item, id).Error; err != nil {
		c.String(http.StatusNotFound, "Not found")
		return
	}

	var updatedItem models.Item
	if err := c.Bind(&updatedItem); err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	// Update fields while preserving ID
	item.Name = updatedItem.Name
	item.Price = updatedItem.Price
	item.Unit = updatedItem.Unit
	item.TaxRate = updatedItem.TaxRate

	h.DB.Save(&item)
	c.HTML(http.StatusOK, "item.html", item)
}

func (h *ItemHandler) ExportItems(c *gin.Context) {
	var items []models.Item
	h.DB.Find(&items)

	csvData, err := csv.ExportItemsToCSV(items)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error exporting items: %v", err)
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=artikli.csv")
	c.Data(http.StatusOK, "text/csv", csvData)
}

func (h *ItemHandler) ImportItems(c *gin.Context) {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Error getting file: %v", err)
		return
	}

	// Create a temporary file to save the uploaded CSV
	tempFile := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveUploadedFile(file, tempFile); err != nil {
		c.String(http.StatusInternalServerError, "Error saving file: %v", err)
		return
	}
	defer os.Remove(tempFile)

	// Drop and recreate the items table to reset IDs
	h.DB.Migrator().DropTable(&models.Item{})
	h.DB.AutoMigrate(&models.Item{})

	// Import items from the CSV file
	if err := csv.ImportItemsFromCSV(tempFile, h.DB); err != nil {
		c.String(http.StatusInternalServerError, "Error importing items: %v", err)
		return
	}

	// Return success response
	c.Redirect(http.StatusFound, "/items")
}
