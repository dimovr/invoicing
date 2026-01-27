# Invoicing App - Rust Implementations

This directory contains four different Rust-based implementations of the invoicing application, each using a different web framework. All implementations provide the same functionality as the original Go application.

## Implementations

| Framework | Directory | Description |
|-----------|-------------|-------------|
| **Actix-Web** | `./actix-web/` | High-performance web framework with actor model |
| **Axum** | `./axum/` | Ergonomic and modular web framework built on Tokio |
| **Rocket** | `./rocket/` | Type-safe web framework with ease of use |
| **Leptos** | `./leptos/` | Modern reactive framework with WASM and SSR |

## Common Features

All implementations include:

- **Company Management**: Create and update company information
- **Items Management**: CRUD operations for inventory items with tax rates
- **Suppliers Management**: Manage supplier information
- **Invoicing**: Create invoices with line items, calculate totals and taxes

## Quick Start

### Actix-Web
```bash
cd rust/actix-web
cargo install diesel_cli --no-default-features --features sqlite
diesel setup
diesel migration run
cargo run
```

### Axum
```bash
cd rust/axum
cargo run
```

### Rocket
```bash
cd rust/rocket
cargo install diesel_cli --no-default-features --features sqlite
diesel setup
diesel migration run
cargo run
```

### Leptos
```bash
cd rust/leptos
cargo install trunk
trunk serve --open
```

## Comparison

| Feature | Actix-Web | Axum | Rocket | Leptos |
|---------|------------|------|--------|---------|
| **Performance** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Ease of Use** | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| **Type Safety** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Async Support** | ✅ | ✅ | ✅ | ✅ |
| **ORM** | Diesel | SQLx | Diesel | SQLx |
| **Templates** | Askama | Custom | Tera | Leptos Components |
| **Frontend** | Server-side | Server-side | Server-side | WASM + SSR |
| **Reactivity** | ❌ | ❌ | ❌ | ✅ |
| **Learning Curve** | Medium | Medium | Low | High |

## Framework Details

### Actix-Web
- **Best for**: High-performance applications, actor-based concurrency
- **Strengths**: Mature ecosystem, excellent performance, flexible
- **Use when**: You need maximum performance and don't mind a steeper learning curve

### Axum
- **Best for**: Modern async applications, Tower ecosystem integration
- **Strengths**: Ergonomic API, async-first, modular design
- **Use when**: You want a modern async framework with great type safety

### Rocket
- **Best for**: Rapid development, ease of use
- **Strengths**: Simple API, type-safe routing, great documentation
- **Use when**: You want to get started quickly with minimal boilerplate

### Leptos
- **Best for**: Modern SPAs with reactivity, full-stack type safety
- **Strengths**: Reactive UI, WASM performance, SSR support
- **Use when**: You want a modern frontend framework with Rust's type safety

## Database Schema

All implementations use the same SQLite database schema:

### Tables
- `company` - Company information
- `item` - Inventory items
- `supplier` - Supplier information
- `invoice` - Invoice headers
- `invoice_item` - Invoice line items

### Relationships
- `invoice.supplier_id` → `supplier.id`
- `invoice_item.invoice_id` → `invoice.id`
- `invoice_item.item_id` → `item.id`

## API Endpoints

All implementations provide the same REST API:

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

## Development

### Running Tests
```bash
# In any implementation directory
cargo test
```

### Building for Release
```bash
# In any implementation directory
cargo build --release
```

### Linting
```bash
cargo clippy
```

### Formatting
```bash
cargo fmt
```

## Choosing the Right Framework

### Choose Actix-Web if:
- You need maximum performance
- You're building a high-traffic API
- You want a mature, battle-tested framework
- You don't mind a steeper learning curve

### Choose Axum if:
- You want a modern async framework
- You're already using the Tower ecosystem
- You prefer ergonomic APIs
- You want compile-time route checking

### Choose Rocket if:
- You want to get started quickly
- You prefer simplicity over performance
- You're new to Rust web development
- You want type-safe routing with minimal boilerplate

### Choose Leptos if:
- You want a reactive frontend
- You need full-stack type safety
- You want SPA-like experience with SSR
- You're willing to learn a new paradigm

## Additional Resources

- [Actix-Web Documentation](https://actix.rs/)
- [Axum Documentation](https://docs.rs/axum/)
- [Rocket Documentation](https://rocket.rs/)
- [Leptos Documentation](https://leptos.dev/)
- [Rust Web Framework Comparison](https://www.arewewebyet.org/)

## License

These are clones of the original Go invoicing application.
