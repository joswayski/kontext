use axum::{routing::get, Router};
use std::time::Duration;
use tower::ServiceBuilder;
use tower_http::compression::CompressionLayer;
use tower_http::timeout::TimeoutLayer;
use tower_http::trace::TraceLayer;

mod handlers;
mod shared;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    tracing_subscriber::fmt::init();

    let app = Router::new()
        .route("/", get(|| async { "Welcome to Kontext API :)" }))
        .fallback(handlers::fallback_404)
        .method_not_allowed_fallback(handlers::fallback_405)
        .layer(
            ServiceBuilder::new()
                // When using a ServiceBuilder, middleware is applied top down
                .layer(TraceLayer::new_for_http())
                .layer(CompressionLayer::new())
                .layer(TimeoutLayer::new(Duration::from_secs(10))),
        );

    let listener = tokio::net::TcpListener::bind("0.0.0.0:4000")
        .await
        .map_err(|e| {
            tracing::error!("Failed to bind to 0.0.0.0:4000: {}", e);
            e
        })?;
    tracing::info!("Server starting on http://0.0.0.0:4000");
    axum::serve(listener, app).await.map_err(|e| {
        tracing::error!("Failed to start server: {}", e);
        e
    })?;
    Ok(())
}
