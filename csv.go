package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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

func WriteProductsToCSV(filename string, products []ProductCsv) error {
	// Create or open the output file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Use a plain text writer instead of csv.Writer to control quoting exactly
	_, err = file.WriteString(formatQuotedLine(headers))
	if err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}

	// Write each product
	for _, product := range products {
		record := make([]string, 14)

		record[0] = fmt.Sprintf("%d", product.ID)               // PLU
		record[1] = product.Name                                // Name
		record[2] = taxRateToVAT(product.TaxRate)               // VAT
		record[3] = "1"                                         // Stock group (default)
		record[4] = fmt.Sprintf("%d", product.PriceType)        // PriceType
		record[5] = fmt.Sprintf("%.2f", float64(product.Price)) // Price
		for i := 6; i < 14; i++ {
			record[i] = "0"
		}

		_, err = file.WriteString(formatQuotedLine(record))
		if err != nil {
			return fmt.Errorf("error writing record: %v", err)
		}
	}

	return nil
}

// formatQuotedLine formats a slice of strings with single quotes and semicolons
func formatQuotedLine(fields []string) string {
	quoted := make([]string, len(fields))
	for i, field := range fields {
		quoted[i] = fmt.Sprintf("\"%s\"", field)
	}
	return strings.Join(quoted, ";") + ";\n"
}

// taxRateToVAT converts tax rate back to VAT code
func taxRateToVAT(rate int) string {
	switch rate {
	case 20:
		return "Ђ"
	case 10:
		return "Е"
	default:
		return "Ђ" // Default to 20% if unknown
	}
}

func main_csv() {
	inputFile := "artikli.csv"
	products, err := ReadProductsFromCSV(inputFile)
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return
	}

	for _, p := range products {
		fmt.Printf("ID: %d, Naziv: %s, PDV: %d, Cena: %d\n",
			p.ID, p.Name, p.TaxRate, p.Price)
	}

	outputFile := "output.csv"
	if err := WriteProductsToCSV(outputFile, products); err != nil {
		fmt.Printf("Error writing CSV: %v\n", err)
		return
	}
	fmt.Println("Successfully wrote to", outputFile)
}
