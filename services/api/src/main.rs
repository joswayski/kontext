use ::shared::config;
use axum::{routing::get, Router};
use handlers::{fallback_404::fallback_404, fallback_405::fallback_405, health::health_check};
use std::time::Duration;
use tokio::signal;
use tower::ServiceBuilder;
use tower_http::compression::CompressionLayer;
use tower_http::timeout::TimeoutLayer;
use tower_http::trace::TraceLayer;
mod handlers;
mod shared;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    config::load_env();

    let port = config::get_port();
    let addr = format!("0.0.0.0:{}", port);

    // Configure tracing subscriber with more detailed logging
    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::DEBUG)
        .with_level(true)
        .with_ansi(true)
        .init();

    let app = Router::new()
        .route(
            "/",
            get(|| async { "Welcome to Kontext API :) - Check the docs for more information!" }),
        )
        .route("/health", get(health_check)) // k8s health check
        .route("/api/health", get(health_check)) // api health check
        .fallback(fallback_404)
        .method_not_allowed_fallback(fallback_405)
        .layer(
            ServiceBuilder::new()
                // When using a ServiceBuilder, middleware is applied top down
                .layer(TraceLayer::new_for_http())
                .layer(CompressionLayer::new())
                .layer(TimeoutLayer::new(Duration::from_secs(10))),
        );

    let listener = tokio::net::TcpListener::bind(&addr).await.map_err(|e| {
        tracing::error!("Failed to bind to {}: {}", addr, e);
        e
    })?;
    tracing::info!("Server starting on http://{}", addr);

    // Create a shutdown signal
    let shutdown = async {
        let ctrl_c = async {
            signal::ctrl_c()
                .await
                .expect("failed to install Ctrl+C handler");
        };

        #[cfg(unix)]
        let terminate = async {
            signal::unix::signal(signal::unix::SignalKind::terminate())
                .expect("failed to install signal handler")
                .recv()
                .await;
        };

        #[cfg(not(unix))]
        let terminate = std::future::pending::<()>();

        tokio::select! {
            _ = ctrl_c => {},
            _ = terminate => {},
        }

        tracing::info!("Shutting down gracefully...");
    };

    // Start the server with graceful shutdown
    axum::serve(listener, app)
        .with_graceful_shutdown(shutdown)
        .await
        .map_err(|e| {
            tracing::error!("Failed to start server: {}", e);
            e
        })?;

    Ok(())
}
