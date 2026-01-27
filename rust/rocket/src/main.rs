mod handlers;
mod models;
mod schema;

use rocket::fairing::{Fairing, Info, Kind};
use rocket::fs::{FileServer, relative};
use rocket::{Build, Rocket};
use rocket::request::{FromRequest, Outcome};
use rocket::http::Status;

pub struct Db;

#[rocket::async_trait]
impl Fairing for Db {
    fn info(&self) -> Info {
        Info {
            name: "Database Fairing",
            kind: Kind::Ignite,
        }
    }

    async fn on_ignite(&self, rocket: Rocket<Build>) -> Outcome<Rocket<Build>, Rocket<Build>> {
        // Run migrations on startup
        // In a real implementation, you would run diesel migrations here
        Outcome::Ok(rocket)
    }
}

#[launch]
async fn rocket() -> _ {
    rocket::build()
        .attach(Db)
        .mount("/", routes![
            handlers::get_invoices,
            // Company routes
            handlers::get_company,
            handlers::upsert_company,
            // Item routes
            handlers::get_items,
            handlers::get_items_partial,
            handlers::get_item_create_form,
            handlers::create_item,
            handlers::get_item_edit_form,
            handlers::update_item,
            handlers::delete_item,
            handlers::export_items,
            // Supplier routes
            handlers::get_suppliers,
            handlers::get_suppliers_partial,
            handlers::get_supplier_create_form,
            handlers::create_supplier,
            handlers::get_supplier_edit_form,
            handlers::update_supplier,
            handlers::delete_supplier,
            // Invoice routes
            handlers::get_invoices,
            handlers::initialize_invoice,
            handlers::get_invoice_edit_page,
            handlers::add_line_item,
            handlers::remove_line_item,
            handlers::complete_invoice,
            handlers::get_invoice_details,
            handlers::delete_invoice,
        ])
        .mount("/static", FileServer::from(relative!("static")))
}
