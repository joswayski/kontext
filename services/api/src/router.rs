use axum::{
    routing::{get, post},
    Router,
};
use std::sync::Arc;
use std::time::Duration;
use tower::ServiceBuilder;
use tower_http::compression::CompressionLayer;
use tower_http::timeout::TimeoutLayer;
use tower_http::trace::TraceLayer;

use crate::handlers::{fallback_404, fallback_405, health_check, root, test_handler};
use shared::{
    clients::{KafkaClient, PostgresClient},
    config::Config,
};

#[derive(Clone)]
struct AppState {
    kafka: Arc<KafkaClient>,
    postgres: Arc<PostgresClient>,
}

impl AppState {
    async fn new(config: &Config) -> Self {
        Self {
            kafka: Arc::new(KafkaClient::new(&config.kafka).await),
            postgres: Arc::new<PostgresClient>(PostgresClient::new(&config.postgres).await),
        }
    }
}

pub async fn create_routes(config: &Config) -> Router {
    let state = AppState::new(&config).await;

    Router::new()
        .route("/", get(root))
        .route("/health", get(health_check)) // k8s health check
        .route("/api/health", get(health_check)) // api health check
        .route("/api/test", post(test_handler))
        .fallback(fallback_404)
        .method_not_allowed_fallback(fallback_405)
        .layer(
            ServiceBuilder::new()
                .layer(TraceLayer::new_for_http())
                .layer(CompressionLayer::new())
                .layer(TimeoutLayer::new(Duration::from_secs(10))),
        )
        .with_state(state)
}
