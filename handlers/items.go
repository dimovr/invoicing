package handlers

import (
	"net/http"
	"strconv"

	"invoicing-item-app/models" // Adjust this import path based on your module name

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
