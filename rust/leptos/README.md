# Invoicing App - Leptos Implementation

A modern, full-stack Rust web application built with Leptos - a reactive framework that compiles to WebAssembly for the frontend and runs on Actix-Web for the backend.

## Features

- **Company Management**: Create and update company information
- **Items Management**: CRUD operations for inventory items with tax rates
- **Suppliers Management**: Manage supplier information
- **Invoicing**: Create invoices with line items, calculate totals and taxes
- **Reactive UI**: Real-time updates with fine-grained reactivity
- **Server-Side Rendering**: Fast initial page loads with SEO benefits
- **Client-Side Hydration**: Interactive SPA-like experience after initial load

## Tech Stack

- **Leptos 0.7**: Modern reactive web framework for Rust
- **Leptos Router**: Client-side routing for SPA navigation
- **Leptos Meta**: Meta tag management for SEO
- **Actix-Web**: High-performance web server (SSR)
- **WASM**: WebAssembly compilation for client-side code
- **SQLx**: Async SQL toolkit for database operations
- **SQLite**: Lightweight database
- **Chrono**: Date and time handling
- **Serde**: Serialization/deserialization

## Project Structure

```
rust/leptos/
├── src/
│   ├── main.rs                   # Application entry point
│   ├── app.rs                    # Root App component
│   ├── components.rs              # Page and UI components
│   └── models.rs                 # Data models
├── public/                       # Static assets
├── Cargo.toml                    # Dependencies
└── README.md                     # This file
```

## Getting Started

### Prerequisites

- Rust 1.70 or later
- Trunk (for WASM development): `cargo install trunk`
- Node.js (for some build tools)

### Development

1. Install Trunk:
   ```bash
   cargo install trunk
   ```

2. Run development server:
   ```bash
   trunk serve --open
   ```

The development server will start on `http://127.0.0.1:8080`

### Production Build

1. Build for production:
   ```bash
   trunk build --release
   ```

2. Run the server:
   ```bash
   cargo run --release --features ssr
   ```

## Features

### Rendering Modes

Leptos supports multiple rendering modes:

- **CSR (Client-Side Rendering)**: Pure client-side rendering
- **SSR (Server-Side Rendering)**: Initial HTML rendered on server, hydrated on client
- **Hydrate**: Hydrate existing HTML on the client

### Component Structure

The application is organized into components:

- `App`: Root component with routing
- `InvoiceList`: List all invoices
- `InvoiceForm`: Create/edit invoices
- `ItemList`: List all items
- `ItemForm`: Create/edit items
- `SupplierList`: List all suppliers
- `SupplierForm`: Create/edit suppliers
- `CompanyPage`: Company information form

### Routing

Uses Leptos Router for client-side navigation:
- `/` - Invoices list
- `/items` - Items list
- `/suppliers` - Suppliers list
- `/company` - Company information

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

## Key Advantages of Leptos

1. **Type Safety**: Full type safety across frontend and backend
2. **Reactivity**: Fine-grained reactive signals for efficient updates
3. **Isomorphic**: Same code runs on server and client
4. **Performance**: WebAssembly for fast client-side execution
5. **SEO**: Server-side rendering for better search engine visibility
6. **No JavaScript**: Write everything in Rust, compile to WASM

## Development

### Running tests
```bash
cargo test
```

### Building for release
```bash
cargo build --release
```

### WASM Development
```bash
# Watch mode for development
trunk serve --open

# Build WASM only
trunk build
```

## Comparison with Other Implementations

| Feature | Leptos | Actix-Web | Axum | Rocket |
|---------|---------|------------|------|--------|
| Frontend | WASM/SSR | Templates | Templates | Templates |
| Reactivity | ✅ Fine-grained | ❌ | ❌ | ❌ |
| Type Safety | ✅ Full | ✅ Backend | ✅ Backend | ✅ Backend |
| SEO | ✅ SSR | ✅ | ✅ | ✅ |
| Learning Curve | High | Medium | Medium | Low |
| Performance | High | High | High | High |

## License

This is a clone of the original Go invoicing application.
