mod handlers;
mod models;
mod schema;

use actix_files as fs;
use actix_web::{middleware, web, App, HttpServer};
use diesel::{r2d2::ConnectionManager, SqliteConnection};
use diesel_migrations::{embed_migrations, EmbeddedMigrations, MigrationHarness};
use r2d2::Pool;
use std::env;

pub type DbPool = Pool<ConnectionManager<SqliteConnection>>;

const MIGRATIONS: EmbeddedMigrations = embed_migrations!("migrations");

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init();

    let database_url = env::var("DATABASE_URL").unwrap_or_else(|_| "invoicing.db".to_string());
    
    // Set up database connection pool
    let manager = ConnectionManager::<SqliteConnection>::new(&database_url);
    let pool = Pool::builder()
        .build(manager)
        .expect("Failed to create pool");

    // Run migrations
    {
        let mut conn = pool.get().expect("Failed to get connection");
        conn.run_pending_migrations(MIGRATIONS)
            .expect("Failed to run migrations");
    }

    let bind_address = "127.0.0.1:8080";
    println!("Starting server at http://{}", bind_address);

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(pool.clone()))
            .wrap(middleware::Logger::default())
            .service(fs::Files::new("/static", "./static").show_files_listing())
            .service(
                web::scope("/")
                    .route("", web::get().to(handlers::get_invoices))
                    // Company routes
                    .route("/company", web::get().to(handlers::get_company))
                    .route("/company", web::post().to(handlers::upsert_company))
                    // Item routes
                    .route("/items", web::get().to(handlers::get_items))
                    .route("/items/list", web::get().to(handlers::get_items_partial))
                    .route("/items/form", web::get().to(handlers::get_item_create_form))
                    .route("/items", web::post().to(handlers::create_item))
                    .route("/items/{id}/edit", web::get().to(handlers::get_item_edit_form))
                    .route("/items/{id}", web::put().to(handlers::update_item))
                    .route("/items/{id}", web::delete().to(handlers::delete_item))
                    .route("/items/export", web::get().to(handlers::export_items))
                    // Supplier routes
                    .route("/suppliers", web::get().to(handlers::get_suppliers))
                    .route("/suppliers/list", web::get().to(handlers::get_suppliers_partial))
                    .route("/suppliers/form", web::get().to(handlers::get_supplier_create_form))
                    .route("/suppliers", web::post().to(handlers::create_supplier))
                    .route("/suppliers/{id}/edit", web::get().to(handlers::get_supplier_edit_form))
                    .route("/suppliers/{id}", web::put().to(handlers::update_supplier))
                    .route("/suppliers/{id}", web::delete().to(handlers::delete_supplier))
                    // Invoice routes
                    .route("/invoices", web::get().to(handlers::get_invoices))
                    .route("/invoices", web::post().to(handlers::initialize_invoice))
                    .route("/invoices/{id}/items", web::post().to(handlers::add_line_item))
                    .route("/invoices/{id}/items/{item_id}", web::delete().to(handlers::remove_line_item))
                    .route("/invoices/{id}/complete", web::post().to(handlers::complete_invoice))
                    .route("/invoices/{id}/view", web::get().to(handlers::get_invoice_details))
                    .route("/invoices/{id}/edit", web::get().to(handlers::get_invoice_edit_page))
                    .route("/invoices/{id}", web::delete().to(handlers::delete_invoice)),
            )
    })
    .bind(bind_address)?
    .run()
    .await
}
