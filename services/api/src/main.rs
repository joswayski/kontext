use ::shared::config::Config;
use api::router::create_routes;
use std::sync::Arc;
use tokio::signal;

#[derive(Clone)]
struct AppState {
    kafka_consumer: Option<Arc<rdkafka::consumer::BaseConsumer>>,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let config = Config::init();

    // Initialize Kafka consumer
    let kafka_consumer = if let Some(cluster) = config.kafka.get_cluster("CLUSTER1") {
        shared::clients::kafka::KafkaClients::create_consumer(cluster.brokers.clone()).map(Arc::new)
    } else {
        None
    };

    let app_state = AppState { kafka_consumer };

    let app = create_routes();

    let addr = format!("0.0.0.0:{}", config.port);
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
