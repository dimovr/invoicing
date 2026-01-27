mod app;
mod components;
mod models;

#[cfg(feature = "hydrate")]
#[wasm_bindgen::prelude::wasm_bindgen(start)]
pub fn hydrate() {
    use app::App;
    leptos::leptos_dom::hydrate_body(App);
}

#[cfg(feature = "csr")]
leptos::mount_to_body(app::App);

#[cfg(feature = "ssr")]
#[actix_web::main]
async fn main() -> std::io::Result<()> {
    use actix_files::Files;
    use actix_web::*;
    use leptos::*;
    use leptos_actix::{generate_route_list, LeptosRoutes};
    use leptos_meta::MetaTags;

    let conf = get_configuration(None).await.unwrap();
    let addr = conf.leptos_options.site_addr.clone();
    let routes = generate_route_list(app::App);

    HttpServer::new(move || {
        let leptos_options = &conf.leptos_options;
        let site_root = &leptos_options.site_root;

        App::new()
            .leptos_routes(leptos_options.to_owned(), routes.to_owned(), app::App)
            .service(Files::new("/", site_root))
            .wrap(middleware::Compress::default())
            .wrap(middleware::Logger::default())
    })
    .bind(&addr)?
    .run()
    .await
}
