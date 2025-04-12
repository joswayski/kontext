use axum::{routing::get, Router};
use std::time::Duration;
use tower::ServiceBuilder;
use tower_http::compression::CompressionLayer;
use tower_http::timeout::TimeoutLayer;
use tower_http::trace::TraceLayer;

use crate::handlers::{fallback_404, fallback_405, health_check, root};

pub fn create_routes() -> Router {
    Router::new()
        .route("/", get(root))
        .route("/health", get(health_check)) // k8s health check
        .route("/api/health", get(health_check)) // api health check
        .fallback(fallback_404)
        .method_not_allowed_fallback(fallback_405)
        .layer(
            ServiceBuilder::new()
                .layer(TraceLayer::new_for_http())
                .layer(CompressionLayer::new())
                .layer(TimeoutLayer::new(Duration::from_secs(10))),
        )
}
