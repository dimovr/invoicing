package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"

    "invoicing-item-app/models"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// ProductCsv represents the simplified product structure
type ProductCsv struct {
    ID        int
    Name      string
    TaxRate   int
    PriceType int
    Price     int
}

// CSV headers as defined in the sample
var headers = []string{
    "PLU", "Name", "VAT", "Stock group", "PriceType", "Price",
    "Single Sale", "Turnover", "Sold Qty", "Stock Qty",
    "Barcode1", "Barcode2", "Barcode3", "Barcode4",
}

// ReadProductsFromCSV reads products from a CSV file
func ReadProductsFromCSV(filename string) ([]ProductCsv, error) {
    // Open the CSV file
    file, err := os.Open(filename)
    if err != nil {
       return nil, fmt.Errorf("error opening file: %v", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = ';'
    reader.LazyQuotes = true

    // Read and skip header row
    _, err = reader.Read()
    if err != nil {
       return nil, fmt.Errorf("error reading header: %v", err)
    }

    var products []ProductCsv

    for {
       record, err := reader.Read()
       if err == io.EOF {
          break
       }
       if err != nil {
          return nil, fmt.Errorf("error reading record: %v", err)
       }

       // Parse the record into ProductCsv
       product, err := parseRecord(record)
       if err != nil {
          return nil, fmt.Errorf("error parsing record: %v", err)
       }

       products = append(products, product)
    }

    return products, nil
}

func parseRecord(record []string) (ProductCsv, error) {
    if len(record) < 14 {
       return ProductCsv{}, fmt.Errorf("invalid record length")
    }

    // Parse PLU (ID)
    id, err := strconv.Atoi(strings.Trim(record[0], "\""))
    if err != nil {
       return ProductCsv{}, fmt.Errorf("invalid PLU: %v", err)
    }

    // Parse Name
    name := strings.Trim(record[1], "\"")

    // Parse VAT (TaxRate)
    vat := strings.Trim(record[2], "\"")
    taxRate := 0
    switch vat {
    case "Ђ":
       taxRate = 20
    case "Е":
       taxRate = 10
    default:
       return ProductCsv{}, fmt.Errorf("invalid VAT code: %s", vat)
    }

    // Parse Price
    priceStr := strings.Trim(record[5], "\"")
    priceFloat, err := strconv.ParseFloat(priceStr, 64)
    if err != nil {
       return ProductCsv{}, fmt.Errorf("invalid price: %v", err)
    }
    price := int(priceFloat)

    return ProductCsv{
       ID:        id,
       Name:      name,
       TaxRate:   taxRate,
       PriceType: 1, // Always set to 1 as specified
       Price:     price,
    }, nil
}

func ConvertToItem(product ProductCsv) models.Item {
    return models.Item{
       Name:    product.Name,
       Price:   float64(product.Price),
       Unit:    "kom", // Default unit, adjust as needed
       TaxRate: product.TaxRate,
    }
}

func main() {
    db, err := gorm.Open(sqlite.Open("invoicing.db"), &gorm.Config{})
    if err != nil {
       panic("failed to connect database")
    }
    db.AutoMigrate(&models.Company{}, &models.Item{}, &models.Supplier{}, &models.InvoiceItem{}, &models.Invoice{})

    inputFile := "artikli.csv"
    products, err := ReadProductsFromCSV(inputFile)
    if err != nil {
       fmt.Printf("Error reading CSV: %v\n", err)
       return
    }

    fmt.Printf("Found %d products to import\n", len(products))

    successCount := 0
    errorCount := 0

    for _, p := range products {
       item := ConvertToItem(p)

       if err := db.Create(&item).Error; err != nil {
          fmt.Printf("Error importing item '%s': %v\n", p.Name, err)
          errorCount++
       } else {
          fmt.Printf("✓ Imported: %s (Price: %.2f, Tax: %d%%)\n",
             item.Name, item.Price, item.TaxRate)
          successCount++
       }
    }

    fmt.Printf("\nImport complete: %d successful, %d failed\n", successCount, errorCount)
}