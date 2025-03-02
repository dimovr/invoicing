package models

import (
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

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
	Address string `json:"address"`
}

type Invoice struct {
	gorm.Model
	ID         uint              `gorm:"primaryKey" json:"id"`
	SupplierID uint              `gorm:"not null" json:"supplier_id"`
	Supplier   Supplier          `gorm:"foreignKey:SupplierID" json:"supplier"`
	LineItems  []InvoiceLineItem `gorm:"foreignKey:InvoiceID" json:"line_items"`
	Subtotal   float64           `json:"subtotal"`   // Calculated: Sum of item prices * quantitys
	TaxAmount  float64           `json:"tax_amount"` // Calculated: Sum of (item price * quantity * tax_rate)
	Total      float64           `json:"total"`      // Calculated: Subtotal + TaxAmount
	Status     string            `gorm:"default: draft" json:"status"`
}

type InvoiceLineItem struct {
	gorm.Model
	// ID      uint    `gorm:"primaryKey" json:"ID"`
	InvoiceID    uint    `gorm:"not null" json:"invoice_id"`
	ItemID       uint    `gorm:"not null" json:"item_id"`
	Item         Item    `gorm:"foreignKey:ItemID" json:"item"`
	Quantity     int     `gorm:"not null" json:"quantity"`
	Subtotal     float64 `json:"subtotal"`   // Calculated: Item.Price * Quantity
	TaxAmount    float64 `json:"tax_amount"` // Calculated: Item.Price * Quantity * (Item.TaxRate / 100)
	Total        float64 `json:"total"`      // Calculated: Subtotal + TaxAmount
	SellingPrice float64 `gorm:"not null" json:"selling_price"`
}
