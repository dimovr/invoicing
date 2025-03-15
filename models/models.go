package models

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	ID      uint    `gorm:"primaryKey" json:"ID"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	TaxRate int     `gorm:"default:0" json:"taxRate"`
	Unit    string  `json:"unit"`
}

type Company struct {
	gorm.Model
	Code       string `gorm:"size:255;not null" json:"code"`
	SectorCode string `gorm:"size:255;not null" json:"sector_code"`
	Sector     string `gorm:"size:255;not null" json:"sector"`
	Name       string `gorm:"size:255;not null" json:"name"`
	Address    string `gorm:"size:255" json:"address"`
	Owner      string `gorm:"size:255" json:"owner"`
	User       string `gorm:"size:255" json:"user"`
}

type Supplier struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey" json:"ID"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Address string `json:"address"`
}

type Invoice struct {
	gorm.Model
	ID             uint          `gorm:"primaryKey" json:"id"`
	SupplierID     uint          `gorm:"not null" json:"supplier_id"`
	Supplier       Supplier      `gorm:"foreignKey:SupplierID" json:"supplier"`
	LineItems      []InvoiceItem `gorm:"foreignKey:InvoiceID" json:"line_items"`
	Subtotal       float64       `json:"subtotal"`   // Calculated: Sum of item prices * quantitys
	TaxAmount      float64       `json:"tax_amount"` // Calculated: Sum of (item price * quantity * tax_rate)
	Total          float64       `json:"total"`      // Calculated: Subtotal + TaxAmount
	Date           time.Time     `json:"date"`
	DocumentNumber string        `json:"document_number"`
}

type InvoiceItem struct {
	ItemID          uint    `json:"item_id"`
	Name            string  `json:"name"`
	Unit            string  `json:"unit"`
	InvoiceID       uint    `json:"invoice_id"`
	Quantity        float64 `json:"quantity"`
	Price           float64 `json:"price"`
	DependentCosts  float64 `json:"dependent_costs"`
	PriceDifference float64 `json:"price_difference"`
	Value           float64 `json:"value"`
	ValueWithoutVat float64 `json:"value_without_vat"`
	VatRate         float64 `json:"vat_rate"`
	VatAmount       float64 `json:"vat_amount"`
	ValueWithVat    float64 `json:"value_with_vat"`
	UnitPrice       float64 `json:"unit_price"`
	Note            string  `json:"note"`
}
