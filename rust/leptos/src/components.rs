use leptos::*;
use leptos_router::*;
use crate::models::*;

// ============== Invoice Components ==============

#[component]
pub fn InvoiceList() -> impl IntoView {
    view! {
        <div class="page">
            <h1>"Invoices"</h1>
            <div class="actions">
                <button class="btn btn-primary">"New Invoice"</button>
            </div>
            <table class="table">
                <thead>
                    <tr>
                        <th>"ID"</th>
                        <th>"Document Number"</th>
                        <th>"Date"</th>
                        <th>"Supplier"</th>
                        <th>"Total"</th>
                        <th>"Actions"</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td>"1"</td>
                        <td>"INV-001"</td>
                        <td>"2024-01-01"</td>
                        <td>"Supplier A"</td>
                        <td>"$100.00"</td>
                        <td>
                            <button class="btn btn-sm">"View"</button>
                            <button class="btn btn-sm btn-danger">"Delete"</button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    }
}

#[component]
pub fn InvoiceForm() -> impl IntoView {
    view! {
        <div class="page">
            <h1>"Edit Invoice"</h1>
            <form>
                <div class="form-group">
                    <label>"Document Number"</label>
                    <input type="text" name="document_number"/>
                </div>
                <div class="form-group">
                    <label>"Date"</label>
                    <input type="date" name="date"/>
                </div>
                <h2>"Line Items"</h2>
                <table class="table">
                    <thead>
                        <tr>
                            <th>"Item"</th>
                            <th>"Quantity"</th>
                            <th>"Price"</th>
                            <th>"Total"</th>
                            <th>"Actions"</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>"Item 1"</td>
                            <td>"10"</td>
                            <td>"$10.00"</td>
                            <td>"$100.00"</td>
                            <td>
                                <button class="btn btn-sm btn-danger">"Remove"</button>
                            </td>
                        </tr>
                    </tbody>
                </table>
                <button type="submit" class="btn btn-success">"Complete Invoice"</button>
            </form>
        </div>
    }
}

// ============== Item Components ==============

#[component]
pub fn ItemList() -> impl IntoView {
    view! {
        <div class="page">
            <h1>"Items"</h1>
            <div class="actions">
                <button class="btn btn-primary">"New Item"</button>
                <button class="btn btn-secondary">"Export CSV"</button>
            </div>
            <table class="table">
                <thead>
                    <tr>
                        <th>"ID"</th>
                        <th>"Name"</th>
                        <th>"Price"</th>
                        <th>"Tax Rate"</th>
                        <th>"Unit"</th>
                        <th>"Actions"</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td>"1"</td>
                        <td>"Widget A"</td>
                        <td>"$10.00"</td>
                        <td>"20%"</td>
                        <td>"pcs"</td>
                        <td>
                            <button class="btn btn-sm">"Edit"</button>
                            <button class="btn btn-sm btn-danger">"Delete"</button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    }
}

#[component]
pub fn ItemForm() -> impl IntoView {
    view! {
        <div class="page">
            <h1>"New Item"</h1>
            <form>
                <div class="form-group">
                    <label>"Name"</label>
                    <input type="text" name="name"/>
                </div>
                <div class="form-group">
                    <label>"Price"</label>
                    <input type="number" step="0.01" name="price"/>
                </div>
                <div class="form-group">
                    <label>"Tax Rate (%)"</label>
                    <input type="number" name="tax_rate"/>
                </div>
                <div class="form-group">
                    <label>"Unit"</label>
                    <input type="text" name="unit"/>
                </div>
                <button type="submit" class="btn btn-primary">"Save"</button>
                <button type="button" class="btn btn-secondary">"Cancel"</button>
            </form>
        </div>
    }
}

// ============== Supplier Components ==============

#[component]
pub fn SupplierList() -> impl IntoView {
    view! {
        <div class="page">
            <h1>"Suppliers"</h1>
            <div class="actions">
                <button class="btn btn-primary">"New Supplier"</button>
            </div>
            <table class="table">
                <thead>
                    <tr>
                        <th>"ID"</th>
                        <th>"Name"</th>
                        <th>"Code"</th>
                        <th>"Address"</th>
                        <th>"Actions"</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td>"1"</td>
                        <td>"Supplier A"</td>
                        <td>"SUP001"</td>
                        <td>"123 Main St"</td>
                        <td>
                            <button class="btn btn-sm">"Edit"</button>
                            <button class="btn btn-sm btn-danger">"Delete"</button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    }
}

#[component]
pub fn SupplierForm() -> impl IntoView {
    view! {
        <div class="page">
            <h1>"New Supplier"</h1>
            <form>
                <div class="form-group">
                    <label>"Name"</label>
                    <input type="text" name="name"/>
                </div>
                <div class="form-group">
                    <label>"Code"</label>
                    <input type="text" name="code"/>
                </div>
                <div class="form-group">
                    <label>"Address"</label>
                    <input type="text" name="address"/>
                </div>
                <button type="submit" class="btn btn-primary">"Save"</button>
                <button type="button" class="btn btn-secondary">"Cancel"</button>
            </form>
        </div>
    }
}

// ============== Company Components ==============

#[component]
pub fn CompanyPage() -> impl IntoView {
    view! {
        <div class="page">
            <h1>"Company Information"</h1>
            <form>
                <div class="form-group">
                    <label>"Code"</label>
                    <input type="text" name="code"/>
                </div>
                <div class="form-group">
                    <label>"Sector Code"</label>
                    <input type="text" name="sector_code"/>
                </div>
                <div class="form-group">
                    <label>"Sector"</label>
                    <input type="text" name="sector"/>
                </div>
                <div class="form-group">
                    <label>"Name"</label>
                    <input type="text" name="name"/>
                </div>
                <div class="form-group">
                    <label>"Address"</label>
                    <input type="text" name="address"/>
                </div>
                <div class="form-group">
                    <label>"Owner"</label>
                    <input type="text" name="owner"/>
                </div>
                <div class="form-group">
                    <label>"User"</label>
                    <input type="text" name="user"/>
                </div>
                <button type="submit" class="btn btn-primary">"Save"</button>
            </form>
        </div>
    }
}
