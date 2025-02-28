package models

import "gorm.io/gorm"

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
	TaxRate int     `gorm:"default:0"`
	Unit    string  `json:"unit"`
}
