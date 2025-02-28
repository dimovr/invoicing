package handlers

import (
	"net/http"
	"todo-item-app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompanyHandler struct {
	DB *gorm.DB
}

func NewCompanyHandler(db *gorm.DB) *CompanyHandler {
	return &CompanyHandler{DB: db}
}

func (h *CompanyHandler) GetCompany(c *gin.Context) {
	var company models.Company
	if err := h.DB.First(&company).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		// No company exists, render empty form
	}

	// If this is an HTMX request, render only the form partial
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(http.StatusOK, "company.html", gin.H{
			"company": company,
		})
		return
	}

	// Otherwise, render the full page
	c.HTML(http.StatusOK, "index.html", gin.H{
		"company": company,
		"active":  "company",
		"Title":   "Company",
	})
}

func (h *CompanyHandler) UpsertCompany(c *gin.Context) {
	var company models.Company
	if err := c.Bind(&company); err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	var existingCompany models.Company
	if err := h.DB.First(&existingCompany).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new company
			if err := h.DB.Create(&company).Error; err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			// Render only the form partial for HTMX
			c.HTML(http.StatusCreated, "company.html", gin.H{
				"company": company,
			})
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Update existing company
	existingCompany.Code = company.Code
	existingCompany.SectorCode = company.SectorCode
	existingCompany.Sector = company.Sector
	existingCompany.Name = company.Name
	existingCompany.Address = company.Address
	existingCompany.Owner = company.Owner
	existingCompany.User = company.User
	if err := h.DB.Save(&existingCompany).Error; err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Render only the form partial for HTMX
	c.HTML(http.StatusOK, "company.html", gin.H{
		"company": existingCompany,
	})
}
