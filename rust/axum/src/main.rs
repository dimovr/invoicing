mod handlers;
mod models;

use axum::{
    routing::{get, post, put, delete},
    Router,
};
use sqlx::SqlitePool;
use std::net::SocketAddr;
use tower_http::services::ServeDir;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| "invoicing_axum=debug,tower_http=debug,axum=trace".into()),
        )
        .with(tracing_subscriber::fmt::layer())
        .init();

    let database_url = std::env::var("DATABASE_URL").unwrap_or_else(|_| "invoicing.db".to_string());

    let pool = SqlitePool::connect(&database_url).await?;
    sqlx::migrate!("./migrations").run(&pool).await?;

    let app = Router::new()
        // Static files
        .nest_service("/static", ServeDir::new("./static"))
        // Company routes
        .route("/company", get(handlers::get_company).post(handlers::upsert_company))
        // Item routes
        .route("/items", get(handlers::get_items).post(handlers::create_item))
        .route("/items/list", get(handlers::get_items_partial))
        .route("/items/form", get(handlers::get_item_create_form))
        .route("/items/:id/edit", get(handlers::get_item_edit_form))
        .route("/items/:id", put(handlers::update_item).delete(handlers::delete_item))
        .route("/items/export", get(handlers::export_items))
        // Supplier routes
        .route("/suppliers", get(handlers::get_suppliers).post(handlers::create_supplier))
        .route("/suppliers/list", get(handlers::get_suppliers_partial))
        .route("/suppliers/form", get(handlers::get_supplier_create_form))
        .route("/suppliers/:id/edit", get(handlers::get_supplier_edit_form))
        .route("/suppliers/:id", put(handlers::update_supplier).delete(handlers::delete_supplier))
        // Invoice routes
        .route("/invoices", get(handlers::get_invoices).post(handlers::initialize_invoice))
        .route("/invoices/:id/items", post(handlers::add_line_item))
        .route("/invoices/:id/items/:item_id", delete(handlers::remove_line_item))
        .route("/invoices/:id/complete", post(handlers::complete_invoice))
        .route("/invoices/:id/view", get(handlers::get_invoice_details))
        .route("/invoices/:id/edit", get(handlers::get_invoice_edit_page))
        .route("/invoices/:id", delete(handlers::delete_invoice))
        // Default route
        .route("/", get(handlers::get_invoices))
        .with_state(pool);

    let addr = SocketAddr::from(([127, 0, 0, 1], 8080));
    tracing::info!("Listening on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await?;
    axum::serve(listener, app).await?;

    Ok(())
}
