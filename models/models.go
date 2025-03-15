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
	Subtotal       float64       `json:"subtotal"`
	TaxAmount      float64       `json:"tax_amount"`
	Total          float64       `json:"total"`
	Date           time.Time     `json:"date"`
	DocumentNumber string        `json:"document_number"`
}

type InvoiceItem struct {
	InvoiceID    uint    `json:"invoice_id" gorm:"uniqueIndex:idx_invoice_item"` // Part of unique constraint
	ItemID       uint    `json:"item_id" gorm:"uniqueIndex:idx_invoice_item"`    // Part of unique constraint
	Name         string  `json:"name"`
	Unit         string  `json:"unit"`
	TaxRate      float64 `json:"tax_rate"`
	Discount     float64 `json:"discount"`
	Quantity     float64 `json:"quantity"`
	BuyingPrice  float64 `json:"buying_price"`
	Subtotal     float64 `json:"subtotal"`
	TaxAmount    float64 `json:"tax_amount"`
	SellingPrice float64 `json:"selling_price"`
	Total        float64 `json:"total"`
	Note         string  `json:"note"`
}
