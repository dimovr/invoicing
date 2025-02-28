package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"invoicing-app/models" // Adjust this import path based on your module name

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
	h.DB.Find(&items)
	fmt.Println(items)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"items":  items,
		"active": "items",
	})
}

func (h *ItemHandler) GetItemsPartial(c *gin.Context) {
	var items []models.Item
	h.DB.Find(&items)
	fmt.Println(items)
	c.HTML(http.StatusOK, "items.html", gin.H{
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
