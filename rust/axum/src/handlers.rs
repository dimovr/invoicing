use axum::{
    extract::{Path, Query, State},
    http::StatusCode,
    response::{Html, IntoResponse, Redirect, Response},
    Form, Json,
};
use chrono::NaiveDate;
use serde::Deserialize;
use sqlx::SqlitePool;

use crate::models::*;

pub async fn get_company(State(pool): State<SqlitePool>) -> Response {
    Html(r#"
        <!DOCTYPE html>
        <html>
        <head><title>Company</title></head>
        <body><h1>Company Management</h1></body>
        </html>
    "#).into_response()
}

pub async fn upsert_company(
    State(pool): State<SqlitePool>,
    Form(form): Form<NewCompany>,
) -> Response {
    Html("<p>Company upserted</p>").into_response()
}

pub async fn get_items(State(pool): State<SqlitePool>) -> Response {
    Html(r#"
        <!DOCTYPE html>
        <html>
        <head><title>Items</title></head>
        <body><h1>Items Management</h1></body>
        </html>
    "#).into_response()
}

pub async fn get_items_partial(State(pool): State<SqlitePool>) -> Response {
    Html("<div>Items list partial</div>").into_response()
}

pub async fn create_item(
    State(pool): State<SqlitePool>,
    Form(form): Form<NewItem>,
) -> Response {
    Html("<div>Item created</div>").into_response()
}

pub async fn delete_item(
    State(pool): State<SqlitePool>,
    Path(id): Path<i32>,
) -> Response {
    Html("").into_response()
}

pub async fn get_item_create_form() -> Response {
    Html("<form>Item create form</form>").into_response()
}

pub async fn get_item_edit_form(
    State(pool): State<SqlitePool>,
    Path(id): Path<i32>,
) -> Response {
    Html("<form>Item edit form</form>").into_response()
}

pub async fn update_item(
    State(pool): State<SqlitePool>,
    Path(id): Path<i32>,
    Form(form): Form<NewItem>,
) -> Response {
    Html("<div>Item updated</div>").into_response()
}

pub async fn export_items(State(pool): State<SqlitePool>) -> Response {
    (
        StatusCode::OK,
        [("content-type", "text/csv")],
        "id,name,price,tax_rate,unit\n",
    )
        .into_response()
}

pub async fn get_suppliers(State(pool): State<SqlitePool>) -> Response {
    Html(r#"
        <!DOCTYPE html>
        <html>
        <head><title>Suppliers</title></head>
        <body><h1>Suppliers Management</h1></body>
        </html>
    "#).into_response()
}

pub async fn get_suppliers_partial(State(pool): State<SqlitePool>) -> Response {
    Html("<div>Suppliers list partial</div>").into_response()
}

pub async fn create_supplier(
    State(pool): State<SqlitePool>,
    Form(form): Form<NewSupplier>,
) -> Response {
    Html("<div>Supplier created</div>").into_response()
}

pub async fn delete_supplier(
    State(pool): State<SqlitePool>,
    Path(id): Path<i32>,
) -> Response {
    Html("").into_response()
}

pub async fn get_supplier_create_form() -> Response {
    Html("<form>Supplier create form</form>").into_response()
}

pub async fn get_supplier_edit_form(
    State(pool): State<SqlitePool>,
    Path(id): Path<i32>,
) -> Response {
    Html("<form>Supplier edit form</form>").into_response()
}

pub async fn update_supplier(
    State(pool): State<SqlitePool>,
    Path(id): Path<i32>,
    Form(form): Form<NewSupplier>,
) -> Response {
    Html("<div>Supplier updated</div>").into_response()
}

pub async fn get_invoices(State(pool): State<SqlitePool>) -> Response {
    Html(r#"
        <!DOCTYPE html>
        <html>
        <head><title>Invoices</title></head>
        <body><h1>Invoices Management</h1></body>
        </html>
    "#).into_response()
}

pub async fn initialize_invoice(
    State(pool): State<SqlitePool>,
    Form(form): Form<NewInvoice>,
) -> Response {
    Redirect::to("/invoices/1/edit").into_response()
}

pub async fn get_invoice_edit_page(
    State(pool): State<SqlitePool>,
    Path(id): Path<i32>,
) -> Response {
    Html("<form>Invoice edit form</form>").into_response()
}

pub async fn add_line_item(
    State(pool): State<SqlitePool>,
    Path(id): Path<i32>,
    Form(form): Form<NewInvoiceItem>,
) -> Response {
    Html("<div>Line item added</div>").into_response()
}

pub async fn remove_line_item(
    State(pool): State<SqlitePool>,
    Path((id, item_id)): Path<(i32, i32)>,
) -> Response {
    Html("").into_response()
}

pub async fn complete_invoice(
    State(pool): State<SqlitePool>,
    Path(id): Path<i32>,
) -> Response {
    Redirect::to("/invoices/1/view").into_response()
}

pub async fn get_invoice_details(
    State(pool): State<SqlitePool>,
    Path(id): Path<i32>,
) -> Response {
    Html("<div>Invoice details</div>").into_response()
}

pub async fn delete_invoice(
    State(pool): State<SqlitePool>,
    Path(id): Path<i32>,
) -> Response {
    Html("").into_response()
}
