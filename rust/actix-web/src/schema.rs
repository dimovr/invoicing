// @generated automatically by Diesel CLI.

diesel::table! {
    company (id) {
        id -> Integer,
        code -> Text,
        sector_code -> Text,
        sector -> Text,
        name -> Text,
        address -> Nullable<Text>,
        owner -> Nullable<Text>,
        user -> Nullable<Text>,
        created_at -> Timestamp,
        updated_at -> Timestamp,
    }
}

diesel::table! {
    invoice (id) {
        id -> Integer,
        supplier_id -> Integer,
        subtotal -> Float,
        tax_amount -> Float,
        total -> Float,
        date -> Timestamp,
        document_number -> Text,
        created_at -> Timestamp,
        updated_at -> Timestamp,
    }
}

diesel::table! {
    invoice_item (invoice_id, item_id) {
        invoice_id -> Integer,
        item_id -> Integer,
        name -> Text,
        unit -> Text,
        tax_rate -> Float,
        discount -> Float,
        quantity -> Float,
        buying_price -> Float,
        subtotal -> Float,
        tax_amount -> Float,
        selling_price -> Float,
        total -> Float,
        note -> Nullable<Text>,
    }
}

diesel::table! {
    item (id) {
        id -> Integer,
        name -> Text,
        price -> Float,
        tax_rate -> Integer,
        unit -> Text,
        created_at -> Timestamp,
        updated_at -> Timestamp,
    }
}

diesel::table! {
    supplier (id) {
        id -> Integer,
        name -> Text,
        code -> Text,
        address -> Nullable<Text>,
        created_at -> Timestamp,
        updated_at -> Timestamp,
    }
}

diesel::joinable!(invoice -> supplier (supplier_id));

diesel::allow_tables_to_appear_in_same_query!(
    company,
    invoice,
    invoice_item,
    item,
    supplier,
);
