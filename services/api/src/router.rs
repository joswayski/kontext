use axum::{
    routing::{get, post},
    Router,
};
use std::collections::HashMap;
use std::{sync::Arc, time::Duration};
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
    kafka: Arc<tokio::sync::Mutex<KafkaClient>>,
    postgres: Arc<tokio::sync::Mutex<PostgresClient>>,
}

impl AppState {
    async fn new(config: Config) -> Self {
        // Initialize Kafka client
        let kafka_client = KafkaClient::new(config.kafka);
        let kafka = Arc::new(tokio::sync::Mutex::new(kafka_client));

        // Initialize Postgres client
        let postgres_client = PostgresClient::new(config.postgres);
        let postgres = Arc::new(tokio::sync::Mutex::new(postgres_client));

        // Initialize connections
        {
            let mut kafka = kafka.lock().await;
            for cluster_name in config.kafka.clusters.keys() {
                if let Err(e) = kafka.create_consumer(cluster_name) {
                    tracing::error!(
                        "Failed to create Kafka consumer for cluster {}: {}",
                        cluster_name,
                        e
                    );
                    panic!(
                        "Failed to create Kafka consumer for cluster {}: {}",
                        cluster_name, e
                    );
                }
                tracing::info!(
                    "Successfully created Kafka consumer for cluster {}",
                    cluster_name
                );
            }
        }

        {
            let mut postgres = postgres.lock().await;
            if let Err(e) = postgres.get_pool().await {
                tracing::error!("Failed to connect to Postgres: {}", e);
                panic!("Failed to connect to Postgres: {}", e);
            }
        }

        Self { kafka, postgres }
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
