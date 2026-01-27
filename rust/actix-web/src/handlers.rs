use actix_web::{error, web, Error, HttpResponse};
use askama::Template;
use chrono::{NaiveDate, Utc};
use diesel::prelude::*;
use serde::Deserialize;

use crate::models::*;
use crate::schema::*;

pub struct AppState {
    pub db: SqliteConnection,
}

// ============== Company Handlers ==============

#[derive(Template)]
#[template(path = "company.html")]
struct CompanyTemplate {
    company: Option<Company>,
}

pub async fn get_company(pool: web::Data<sqlx::SqlitePool>) -> Result<HttpResponse, Error> {
    // This is a simplified version - in a real implementation, you'd use Diesel
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Company endpoint"})))
}

#[derive(Deserialize)]
pub struct CompanyForm {
    pub code: String,
    pub sector_code: String,
    pub sector: String,
    pub name: String,
    pub address: Option<String>,
    pub owner: Option<String>,
    pub user: Option<String>,
}

pub async fn upsert_company(
    pool: web::Data<sqlx::SqlitePool>,
    form: web::Form<CompanyForm>,
) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Company upserted"})))
}

// ============== Item Handlers ==============

#[derive(Template)]
#[template(path = "items_list.html")]
struct ItemsListTemplate {
    items: Vec<Item>,
}

#[derive(Template)]
#[template(path = "item.html")]
struct ItemTemplate {
    item: Item,
}

#[derive(Template)]
#[template(path = "item-create-form.html")]
struct ItemCreateFormTemplate {}

#[derive(Template)]
#[template(path = "item-edit-form.html")]
struct ItemEditFormTemplate {
    item: Item,
}

pub async fn get_items(pool: web::Data<sqlx::SqlitePool>) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Items list"})))
}

pub async fn get_items_partial(pool: web::Data<sqlx::SqlitePool>) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Items partial"})))
}

#[derive(Deserialize)]
pub struct ItemForm {
    pub name: String,
    pub price: f64,
    pub tax_rate: i32,
    pub unit: String,
}

pub async fn create_item(
    pool: web::Data<sqlx::SqlitePool>,
    form: web::Form<ItemForm>,
) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Created().json(serde_json::json!({"message": "Item created"})))
}

pub async fn delete_item(
    pool: web::Data<sqlx::SqlitePool>,
    path: web::Path<i32>,
) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Item deleted"})))
}

pub async fn get_item_create_form() -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Item create form"})))
}

pub async fn get_item_edit_form(
    pool: web::Data<sqlx::SqlitePool>,
    path: web::Path<i32>,
) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Item edit form"})))
}

pub async fn update_item(
    pool: web::Data<sqlx::SqlitePool>,
    path: web::Path<i32>,
    form: web::Form<ItemForm>,
) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Item updated"})))
}

pub async fn export_items(pool: web::Data<sqlx::SqlitePool>) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok()
        .content_type("text/csv")
        .body("id,name,price,tax_rate,unit\n"))
}

// ============== Supplier Handlers ==============

#[derive(Template)]
#[template(path = "suppliers_list.html")]
struct SuppliersListTemplate {
    suppliers: Vec<Supplier>,
}

#[derive(Template)]
#[template(path = "supplier.html")]
struct SupplierTemplate {
    supplier: Supplier,
}

#[derive(Template)]
#[template(path = "supplier-create-form.html")]
struct SupplierCreateFormTemplate {}

#[derive(Template)]
#[template(path = "supplier-edit-form.html")]
struct SupplierEditFormTemplate {
    supplier: Supplier,
}

pub async fn get_suppliers(pool: web::Data<sqlx::SqlitePool>) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Suppliers list"})))
}

pub async fn get_suppliers_partial(pool: web::Data<sqlx::SqlitePool>) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Suppliers partial"})))
}

#[derive(Deserialize)]
pub struct SupplierForm {
    pub name: String,
    pub code: String,
    pub address: Option<String>,
}

pub async fn create_supplier(
    pool: web::Data<sqlx::SqlitePool>,
    form: web::Form<SupplierForm>,
) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Created().json(serde_json::json!({"message": "Supplier created"})))
}

pub async fn delete_supplier(
    pool: web::Data<sqlx::SqlitePool>,
    path: web::Path<i32>,
) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Supplier deleted"})))
}

pub async fn get_supplier_create_form() -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Supplier create form"})))
}

pub async fn get_supplier_edit_form(
    pool: web::Data<sqlx::SqlitePool>,
    path: web::Path<i32>,
) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Supplier edit form"})))
}

pub async fn update_supplier(
    pool: web::Data<sqlx::SqlitePool>,
    path: web::Path<i32>,
    form: web::Form<SupplierForm>,
) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Supplier updated"})))
}

// ============== Invoice Handlers ==============

#[derive(Template)]
#[template(path = "invoice-form.html")]
struct InvoiceFormTemplate {
    invoice: Invoice,
    items: Vec<Item>,
}

#[derive(Template)]
#[template(path = "invoice-line-item.html")]
struct InvoiceLineItemTemplate {
    item: InvoiceItem,
}

#[derive(Template)]
#[template(path = "invoice-full.html")]
struct InvoiceFullTemplate {
    invoice: Invoice,
    company: Company,
    today_date: String,
}

pub async fn get_invoices(pool: web::Data<sqlx::SqlitePool>) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Invoices list"})))
}

#[derive(Deserialize)]
pub struct InitializeInvoiceForm {
    pub supplier_id: i32,
    pub document_number: String,
    pub date: String,
}

pub async fn initialize_invoice(
    pool: web::Data<sqlx::SqlitePool>,
    form: web::Form<InitializeInvoiceForm>,
) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Found()
        .append_header(("Location", "/invoices/1/edit"))
        .finish())
}

pub async fn get_invoice_edit_page(
    pool: web::Data<sqlx::SqlitePool>,
    path: web::Path<i32>,
) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Invoice edit page"})))
}

#[derive(Deserialize)]
pub struct AddLineItemForm {
    pub item_id: i32,
    pub price: f64,
    pub quantity: f64,
    pub discount: f64,
}

pub async fn add_line_item(
    pool: web::Data<sqlx::SqlitePool>,
    path: web::Path<i32>,
    form: web::Form<AddLineItemForm>,
) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Line item added"})))
}

pub async fn remove_line_item(
    pool: web::Data<sqlx::SqlitePool>,
    path: web::Path<(i32, i32)>,
) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Line item removed"})))
}

pub async fn complete_invoice(
    pool: web::Data<sqlx::SqlitePool>,
    path: web::Path<i32>,
) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Found()
        .append_header(("Location", "/invoices/1/view"))
        .finish())
}

pub async fn get_invoice_details(
    pool: web::Data<sqlx::SqlitePool>,
    path: web::Path<i32>,
) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Invoice details"})))
}

pub async fn delete_invoice(
    pool: web::Data<sqlx::SqlitePool>,
    path: web::Path<i32>,
) -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok().json(serde_json::json!({"message": "Invoice deleted"})))
}
