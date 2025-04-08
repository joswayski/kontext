use ::shared::config;
use api::router::create_routes;
use tokio::signal;

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

    let app = create_routes();

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
