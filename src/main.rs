use axum::{
    http::{Request, Response},
    routing::any,
    Extension, Router,
};
use hyper::{client::HttpConnector, Body, Client, Uri};
use std::{env, net::SocketAddr};
use tower_http::trace::TraceLayer;
use tracing::info;

type HyperClient = Client<hyper_tls::HttpsConnector<HttpConnector>>;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    dotenvy::dotenv().ok();
    if std::env::var_os("RUST_LOG").is_none() {
        std::env::set_var("RUST_LOG", "info");
    }
    // initialize tracing
    tracing_subscriber::fmt::init();

    let https = hyper_tls::HttpsConnector::new();
    let client = Client::builder().build::<_, hyper::Body>(https);
    // build our application with a route
    let app = Router::new()
        .route("/*uri", any(handler))
        .layer(Extension(client))
        .layer(TraceLayer::new_for_http());
    // run our app with hyper
    // `axum::Server` is a re-export of `hyper::Server`
    let addr = SocketAddr::from((
        [127, 0, 0, 1],
        env::var("PORT").unwrap_or("3001".to_owned()).parse()?,
    ));
    info!("listening on {}", addr);
    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await?;
    Ok(())
}

async fn handler(
    Extension(client): Extension<HyperClient>,
    // NOTE: Make sure to put the request extractor last because once the request
    // is extracted, extensions can't be extracted anymore.
    mut req: Request<Body>,
) -> Response<Body> {
    let path = req.uri().path();
    if path == "/" {
        return Response::builder()
            .status(302)
            .header("Location", "https://github.com/Xhofe")
            .body(Body::empty())
            .unwrap();
    }
    let path_query = req
        .uri()
        .path_and_query()
        .map(|v| v.as_str())
        .unwrap_or(path);

    let uri = format!("{}", &path_query[1..]);
    info!("uri: {}", uri);

    *req.uri_mut() = Uri::try_from(uri).unwrap();
    let has_host = req.headers().contains_key("host");
    if has_host {
        req.headers_mut().remove("host");
        let host = req.uri().host().unwrap_or_default().to_owned();
        req.headers_mut().insert("host", host.parse().unwrap());
    }

    client.request(req).await.unwrap()
}
