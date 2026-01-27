use chrono::{NaiveDateTime, Utc};
use diesel::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Debug, Queryable, Serialize, Deserialize, Clone)]
pub struct Company {
    pub id: i32,
    pub code: String,
    pub sector_code: String,
    pub sector: String,
    pub name: String,
    pub address: Option<String>,
    pub owner: Option<String>,
    pub user: Option<String>,
    pub created_at: NaiveDateTime,
    pub updated_at: NaiveDateTime,
}

#[derive(Debug, Insertable, Deserialize, AsChangeset)]
#[diesel(table_name = crate::schema::company)]
pub struct NewCompany {
    pub code: String,
    pub sector_code: String,
    pub sector: String,
    pub name: String,
    pub address: Option<String>,
    pub owner: Option<String>,
    pub user: Option<String>,
}

#[derive(Debug, Queryable, Serialize, Deserialize, Clone)]
pub struct Item {
    pub id: i32,
    pub name: String,
    pub price: f64,
    pub tax_rate: i32,
    pub unit: String,
    pub created_at: NaiveDateTime,
    pub updated_at: NaiveDateTime,
}

#[derive(Debug, Insertable, Deserialize, AsChangeset)]
#[diesel(table_name = crate::schema::item)]
pub struct NewItem {
    pub name: String,
    pub price: f64,
    pub tax_rate: i32,
    pub unit: String,
}

#[derive(Debug, Queryable, Serialize, Deserialize, Clone)]
pub struct Supplier {
    pub id: i32,
    pub name: String,
    pub code: String,
    pub address: Option<String>,
    pub created_at: NaiveDateTime,
    pub updated_at: NaiveDateTime,
}

#[derive(Debug, Insertable, Deserialize, AsChangeset)]
#[diesel(table_name = crate::schema::supplier)]
pub struct NewSupplier {
    pub name: String,
    pub code: String,
    pub address: Option<String>,
}

#[derive(Debug, Queryable, Serialize, Deserialize, Clone)]
pub struct Invoice {
    pub id: i32,
    pub supplier_id: i32,
    pub subtotal: f64,
    pub tax_amount: f64,
    pub total: f64,
    pub date: NaiveDateTime,
    pub document_number: String,
    pub created_at: NaiveDateTime,
    pub updated_at: NaiveDateTime,
}

#[derive(Debug, Insertable, Deserialize)]
#[diesel(table_name = crate::schema::invoice)]
pub struct NewInvoice {
    pub supplier_id: i32,
    pub subtotal: f64,
    pub tax_amount: f64,
    pub total: f64,
    pub date: NaiveDateTime,
    pub document_number: String,
}

#[derive(Debug, Queryable, Serialize, Deserialize, Clone)]
pub struct InvoiceItem {
    pub invoice_id: i32,
    pub item_id: i32,
    pub name: String,
    pub unit: String,
    pub tax_rate: f64,
    pub discount: f64,
    pub quantity: f64,
    pub buying_price: f64,
    pub subtotal: f64,
    pub tax_amount: f64,
    pub selling_price: f64,
    pub total: f64,
    pub note: Option<String>,
}

#[derive(Debug, Insertable, Deserialize)]
#[diesel(table_name = crate::schema::invoice_item)]
pub struct NewInvoiceItem {
    pub invoice_id: i32,
    pub item_id: i32,
    pub name: String,
    pub unit: String,
    pub tax_rate: f64,
    pub discount: f64,
    pub quantity: f64,
    pub buying_price: f64,
    pub subtotal: f64,
    pub tax_amount: f64,
    pub selling_price: f64,
    pub total: f64,
    pub note: Option<String>,
}

#[derive(Debug, Serialize)]
pub struct InvoiceWithDetails {
    pub id: i32,
    pub supplier_id: i32,
    pub supplier: Option<Supplier>,
    pub line_items: Vec<InvoiceItem>,
    pub subtotal: f64,
    pub tax_amount: f64,
    pub total: f64,
    pub date: NaiveDateTime,
    pub document_number: String,
    pub created_at: NaiveDateTime,
    pub updated_at: NaiveDateTime,
}
