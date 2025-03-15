package handlers

import (
	"invoicing-item-app/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SupplierHandler struct {
	DB *gorm.DB
}

func NewSupplierHandler(db *gorm.DB) *SupplierHandler {
	return &SupplierHandler{DB: db}
}

func (h *SupplierHandler) GetSuppliers(c *gin.Context) {
	var suppliers []models.Supplier
	h.DB.Find(&suppliers)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"suppliers": suppliers,
		"active":    "suppliers",
		"Title":     "Suppliers",
	})
}

func (h *SupplierHandler) GetSuppliersPartial(c *gin.Context) {
	var suppliers []models.Supplier
	h.DB.Find(&suppliers)
	c.HTML(http.StatusOK, "suppliers.html", gin.H{
		"suppliers": suppliers,
	})
}

func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var supplier models.Supplier
	if err := c.Bind(&supplier); err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	h.DB.Create(&supplier)
	c.HTML(http.StatusCreated, "supplier.html", supplier)
}

func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var supplier models.Supplier
	if err := h.DB.First(&supplier, id).Error; err != nil {
		c.String(http.StatusNotFound, "Not found")
		return
	}
	h.DB.Delete(&supplier)
	c.String(http.StatusOK, "")
}

func (h *SupplierHandler) GetSupplierCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "supplier-create-form.html", gin.H{})
}

func (h *SupplierHandler) GetSupplierEditForm(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var supplier models.Supplier
	if err := h.DB.First(&supplier, id).Error; err != nil {
		c.String(http.StatusNotFound, "Not found")
		return
	}
	c.HTML(http.StatusOK, "supplier-edit-form.html", gin.H{"supplier": supplier})
}

func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var supplier models.Supplier
	if err := h.DB.First(&supplier, id).Error; err != nil {
		c.String(http.StatusNotFound, "Not found")
		return
	}

	var updatedSupplier models.Supplier
	if err := c.Bind(&updatedSupplier); err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	// Update fields while preserving ID
	supplier.Name = updatedSupplier.Name
	supplier.Code = updatedSupplier.Code
	supplier.Address = updatedSupplier.Address

	h.DB.Save(&supplier)
	c.HTML(http.StatusOK, "supplier.html", supplier)
}
