use chrono::{NaiveDateTime, Utc};
use diesel::prelude::*;
use serde::{Deserialize, Serialize};
use utoipa::ToSchema;

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

impl NewCompany {
    pub fn new(
        code: String,
        sector_code: String,
        sector: String,
        name: String,
        address: Option<String>,
        owner: Option<String>,
        user: Option<String>,
    ) -> Self {
        Self {
            code,
            sector_code,
            sector,
            name,
            address,
            owner,
            user,
        }
    }
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

impl NewItem {
    pub fn new(name: String, price: f64, tax_rate: i32, unit: String) -> Self {
        Self {
            name,
            price,
            tax_rate,
            unit,
        }
    }
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

impl NewSupplier {
    pub fn new(name: String, code: String, address: Option<String>) -> Self {
        Self {
            name,
            code,
            address,
        }
    }
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

impl NewInvoice {
    pub fn new(supplier_id: i32, document_number: String, date: NaiveDateTime) -> Self {
        Self {
            supplier_id,
            subtotal: 0.0,
            tax_amount: 0.0,
            total: 0.0,
            date,
            document_number,
        }
    }
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

impl NewInvoiceItem {
    pub fn new(
        invoice_id: i32,
        item_id: i32,
        name: String,
        unit: String,
        tax_rate: f64,
        discount: f64,
        quantity: f64,
        buying_price: f64,
        subtotal: f64,
        tax_amount: f64,
        selling_price: f64,
        total: f64,
        note: Option<String>,
    ) -> Self {
        Self {
            invoice_id,
            item_id,
            name,
            unit,
            tax_rate,
            discount,
            quantity,
            buying_price,
            subtotal,
            tax_amount,
            selling_price,
            total,
            note,
        }
    }
}

#[derive(Debug, Serialize, Deserialize)]
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
