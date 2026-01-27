use leptos::*;
use leptos_router::*;

use crate::components::*;

#[component]
pub fn App() -> impl IntoView {
    provide_meta_context();

    view! {
        <html lang="en">
            <head>
                <meta charset="UTF-8"/>
                <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
                <title>"Invoicing App"</title>
                <script src="https://unpkg.com/htmx.org@1.9.10"></script>
            </head>
            <body>
                <div class="container">
                    <nav class="navbar">
                        <ul>
                            <li>
                                <A href="/">"Invoices"</A>
                            </li>
                            <li>
                                <A href="/items">"Items"</A>
                            </li>
                            <li>
                                <A href="/suppliers">"Suppliers"</A>
                            </li>
                            <li>
                                <A href="/company">"Company"</A>
                            </li>
                        </ul>
                    </nav>
                    <main>
                        <Routes>
                            <Route path="/" view=InvoiceList/>
                            <Route path="/items" view=ItemList/>
                            <Route path="/suppliers" view=SupplierList/>
                            <Route path="/company" view=CompanyPage/>
                        </Routes>
                    </main>
                </div>
            </body>
        </html>
    }
}
