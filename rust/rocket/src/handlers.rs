use rocket::response::{content::RawHtml, Redirect};
use rocket::serde::{Deserialize, Serialize};
use rocket::FromForm;

use crate::models::*;

// ============== Company Handlers ==============

#[get("/company")]
pub fn get_company() -> RawHtml<&'static str> {
    RawHtml(r#"
        <!DOCTYPE html>
        <html>
        <head><title>Company</title></head>
        <body><h1>Company Management</h1></body>
        </html>
    "#)
}

#[post("/company", data = "<form>")]
pub fn upsert_company(form: Form<NewCompany>) -> RawHtml<&'static str> {
    RawHtml("<p>Company upserted</p>")
}

// ============== Item Handlers ==============

#[get("/items")]
pub fn get_items() -> RawHtml<&'static str> {
    RawHtml(r#"
        <!DOCTYPE html>
        <html>
        <head><title>Items</title></head>
        <body><h1>Items Management</h1></body>
        </html>
    "#)
}

#[get("/items/list")]
pub fn get_items_partial() -> RawHtml<&'static str> {
    RawHtml("<div>Items list partial</div>")
}

#[get("/items/form")]
pub fn get_item_create_form() -> RawHtml<&'static str> {
    RawHtml("<form>Item create form</form>")
}

#[post("/items", data = "<form>")]
pub fn create_item(form: Form<NewItem>) -> RawHtml<&'static str> {
    RawHtml("<div>Item created</div>")
}

#[delete("/items/<id>")]
pub fn delete_item(id: i32) -> RawHtml<&'static str> {
    RawHtml("")
}

#[get("/items/<id>/edit")]
pub fn get_item_edit_form(id: i32) -> RawHtml<&'static str> {
    RawHtml("<form>Item edit form</form>")
}

#[put("/items/<id>", data = "<form>")]
pub fn update_item(id: i32, form: Form<NewItem>) -> RawHtml<&'static str> {
    RawHtml("<div>Item updated</div>")
}

#[get("/items/export")]
pub fn export_items() -> (rocket::http::ContentType, &'static str) {
    (rocket::http::ContentType::CSV, "id,name,price,tax_rate,unit\n")
}

// ============== Supplier Handlers ==============

#[get("/suppliers")]
pub fn get_suppliers() -> RawHtml<&'static str> {
    RawHtml(r#"
        <!DOCTYPE html>
        <html>
        <head><title>Suppliers</title></head>
        <body><h1>Suppliers Management</h1></body>
        </html>
    "#)
}

#[get("/suppliers/list")]
pub fn get_suppliers_partial() -> RawHtml<&'static str> {
    RawHtml("<div>Suppliers list partial</div>")
}

#[get("/suppliers/form")]
pub fn get_supplier_create_form() -> RawHtml<&'static str> {
    RawHtml("<form>Supplier create form</form>")
}

#[post("/suppliers", data = "<form>")]
pub fn create_supplier(form: Form<NewSupplier>) -> RawHtml<&'static str> {
    RawHtml("<div>Supplier created</div>")
}

#[delete("/suppliers/<id>")]
pub fn delete_supplier(id: i32) -> RawHtml<&'static str> {
    RawHtml("")
}

#[get("/suppliers/<id>/edit")]
pub fn get_supplier_edit_form(id: i32) -> RawHtml<&'static str> {
    RawHtml("<form>Supplier edit form</form>")
}

#[put("/suppliers/<id>", data = "<form>")]
pub fn update_supplier(id: i32, form: Form<NewSupplier>) -> RawHtml<&'static str> {
    RawHtml("<div>Supplier updated</div>")
}

// ============== Invoice Handlers ==============

#[get("/invoices")]
pub fn get_invoices() -> RawHtml<&'static str> {
    RawHtml(r#"
        <!DOCTYPE html>
        <html>
        <head><title>Invoices</title></head>
        <body><h1>Invoices Management</h1></body>
        </html>
    "#)
}

#[derive(FromForm, Deserialize, Serialize)]
pub struct InitializeInvoiceForm {
    pub supplier_id: i32,
    pub document_number: String,
    pub date: String,
}

#[post("/invoices", data = "<form>")]
pub fn initialize_invoice(form: Form<InitializeInvoiceForm>) -> Redirect {
    Redirect::to(uri!("/invoices/1/edit"))
}

#[get("/invoices/<id>/edit")]
pub fn get_invoice_edit_page(id: i32) -> RawHtml<&'static str> {
    RawHtml("<form>Invoice edit form</form>")
}

#[derive(FromForm, Deserialize, Serialize)]
pub struct AddLineItemForm {
    pub item_id: i32,
    pub price: f64,
    pub quantity: f64,
    pub discount: f64,
}

#[post("/invoices/<id>/items", data = "<form>")]
pub fn add_line_item(id: i32, form: Form<AddLineItemForm>) -> RawHtml<&'static str> {
    RawHtml("<div>Line item added</div>")
}

#[delete("/invoices/<id>/items/<item_id>")]
pub fn remove_line_item(id: i32, item_id: i32) -> RawHtml<&'static str> {
    RawHtml("")
}

#[post("/invoices/<id>/complete")]
pub fn complete_invoice(id: i32) -> Redirect {
    Redirect::to(uri!("/invoices/1/view"))
}

#[get("/invoices/<id>/view")]
pub fn get_invoice_details(id: i32) -> RawHtml<&'static str> {
    RawHtml("<div>Invoice details</div>")
}

#[delete("/invoices/<id>")]
pub fn delete_invoice(id: i32) -> RawHtml<&'static str> {
    RawHtml("")
}
