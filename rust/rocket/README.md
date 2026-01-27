# Invoicing App - Rocket Implementation

A Rust-based clone of the Go invoicing application using the Rocket web framework.

## Features

- **Company Management**: Create and update company information
- **Items Management**: CRUD operations for inventory items with tax rates
- **Suppliers Management**: Manage supplier information
- **Invoicing**: Create invoices with line items, calculate totals and taxes

## Tech Stack

- **Rocket 0.5**: Web framework with type safety and ease of use
- **Diesel 2.2**: ORM and query builder for database operations
- **SQLite**: Lightweight database
- **Rocket Dyn Templates**: Template engine for server-side rendering
- **Chrono**: Date and time handling
- **Serde**: Serialization/deserialization

## Project Structure

```
rust/rocket/
├── migrations/                    # Database migrations
│   └── 2023-01-01-000000_initial/
│       ├── up.sql
│       └── down.sql
├── src/
│   ├── main.rs                   # Application entry point and routing
│   ├── models.rs                 # Data models
│   ├── schema.rs                 # Diesel schema definitions
│   └── handlers.rs               # Request handlers
├── Cargo.toml                    # Dependencies
└── README.md                     # This file
```

## Getting Started

### Prerequisites

- Rust 1.70 or later
- Diesel CLI (`cargo install diesel_cli`)

### Setup

1. Install Diesel CLI:
   ```bash
   cargo install diesel_cli --no-default-features --features sqlite
   ```

2. Set up the database:
   ```bash
   diesel setup
   diesel migration run
   ```

3. Run the application:
   ```bash
   cargo run
   ```

The server will start on `http://127.0.0.1:8000`

## API Endpoints

### Company
- `GET /company` - Get company information
- `POST /company` - Create or update company

### Items
- `GET /items` - List all items
- `GET /items/list` - Get items list (partial)
- `GET /items/form` - Get item creation form
- `POST /items` - Create new item
- `GET /items/{id}/edit` - Get item edit form
- `PUT /items/{id}` - Update item
- `DELETE /items/{id}` - Delete item
- `GET /items/export` - Export items as CSV

### Suppliers
- `GET /suppliers` - List all suppliers
- `GET /suppliers/list` - Get suppliers list (partial)
- `GET /suppliers/form` - Get supplier creation form
- `POST /suppliers` - Create new supplier
- `GET /suppliers/{id}/edit` - Get supplier edit form
- `PUT /suppliers/{id}` - Update supplier
- `DELETE /suppliers/{id}` - Delete supplier

### Invoices
- `GET /invoices` - List all invoices
- `POST /invoices` - Initialize new invoice
- `POST /invoices/{id}/items` - Add line item to invoice
- `DELETE /invoices/{id}/items/{item_id}` - Remove line item
- `POST /invoices/{id}/complete` - Complete invoice
- `GET /invoices/{id}/view` - View invoice details
- `GET /invoices/{id}/edit` - Edit invoice
- `DELETE /invoices/{id}` - Delete invoice

## Database Schema

### Company
- `id` (INTEGER, PRIMARY KEY)
- `code` (TEXT, NOT NULL)
- `sector_code` (TEXT, NOT NULL)
- `sector` (TEXT, NOT NULL)
- `name` (TEXT, NOT NULL)
- `address` (TEXT, OPTIONAL)
- `owner` (TEXT, OPTIONAL)
- `user` (TEXT, OPTIONAL)

### Item
- `id` (INTEGER, PRIMARY KEY)
- `name` (TEXT, NOT NULL)
- `price` (REAL, NOT NULL)
- `tax_rate` (INTEGER, NOT NULL, DEFAULT 0)
- `unit` (TEXT, NOT NULL)

### Supplier
- `id` (INTEGER, PRIMARY KEY)
- `name` (TEXT, NOT NULL)
- `code` (TEXT, NOT NULL)
- `address` (TEXT, OPTIONAL)

### Invoice
- `id` (INTEGER, PRIMARY KEY)
- `supplier_id` (INTEGER, NOT NULL, FOREIGN KEY)
- `subtotal` (REAL, NOT NULL, DEFAULT 0)
- `tax_amount` (REAL, NOT NULL, DEFAULT 0)
- `total` (REAL, NOT NULL, DEFAULT 0)
- `date` (TEXT, NOT NULL)
- `document_number` (TEXT, NOT NULL)

### InvoiceItem
- `invoice_id` (INTEGER, PRIMARY KEY, FOREIGN KEY)
- `item_id` (INTEGER, PRIMARY KEY, FOREIGN KEY)
- `name` (TEXT, NOT NULL)
- `unit` (TEXT, NOT NULL)
- `tax_rate` (REAL, NOT NULL)
- `discount` (REAL, NOT NULL)
- `quantity` (REAL, NOT NULL)
- `buying_price` (REAL, NOT NULL)
- `subtotal` (REAL, NOT NULL)
- `tax_amount` (REAL, NOT NULL)
- `selling_price` (REAL, NOT NULL)
- `total` (REAL, NOT NULL)
- `note` (TEXT, OPTIONAL)

## Development

### Running tests
```bash
cargo test
```

### Building for release
```bash
cargo build --release
```

## Key Features of Rocket

1. **Type Safety**: Rocket uses Rust's type system to ensure correctness
2. **Procedural Macros**: Simplifies route definitions and request handling
3. **Fairings**: Middleware system for cross-cutting concerns
4. **State Management**: Built-in managed state for database connections
5. **Form Handling**: Automatic form parsing and validation
6. **Testing**: Built-in testing support for routes

## License

This is a clone of the original Go invoicing application.
