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

// Supplier represents a vendor that provides goods or services
type Supplier struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey" json:"ID"`
	Name    string `gorm:"size:255;not null"  json:"name"`
	Address string `gorm:"size:255;not null" json:"address"`
}

// // Invoice represents a bill from a supplier
// type Invoice struct {
// 	gorm.Model
// 	InvoiceNumber string    `gorm:"size:50;not null;uniqueIndex"`
// 	SupplierID    uint      `gorm:"not null"`
// 	Supplier      Supplier  `gorm:"foreignKey:SupplierID"`
// 	IssueDate     time.Time `gorm:"not null"`
// 	DueDate       time.Time `gorm:"not null"`
// 	PaidDate      *time.Time
// 	SubTotal      float64 `gorm:"type:decimal(10,2);not null"`
// 	TaxAmount     float64 `gorm:"type:decimal(10,2);not null"`
// 	TotalAmount   float64 `gorm:"type:decimal(10,2);not null"`
// 	Status        string  `gorm:"size:20;default:'draft'"` // draft, pending, paid, overdue, cancelled
// 	Notes         string  `gorm:"size:2550"`
// 	AttachmentURL string  `gorm:"size:255"`
// 	Items         []InvoiceItem
// }

// // InvoiceItem represents a line item in an invoice
// type InvoiceItem struct {
// 	gorm.Model
// 	InvoiceID   uint    `gorm:"not null"`
// 	ItemID      uint    `gorm:"not null"`
// 	Item        Item    `gorm:"foreignKey:ItemID"`
// 	Description string  `gorm:"size:500"`
// 	Quantity    float64 `gorm:"type:decimal(10,2);not null"`
// 	UnitPrice   float64 `gorm:"type:decimal(10,2);not null"`
// 	TaxRate     float64 `gorm:"type:decimal(5,2);default:0"`
// 	TaxAmount   float64 `gorm:"type:decimal(10,2);not null"`
// 	TotalAmount float64 `gorm:"type:decimal(10,2);not null"`
// }
