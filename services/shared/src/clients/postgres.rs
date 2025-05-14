use sqlx::postgres::PgPoolOptions;
use sqlx::PgPool;
use tracing;

use crate::config::postgres::PostgresConfig;

pub struct PostgresClient {
    pool: PgPool,
}

impl PostgresClient {
    pub async fn new(config: &PostgresConfig) -> Self {
        let pool = match PgPoolOptions::new()
            .max_connections(config.max_connections)
            .connect(&config.url)
            .await
        {
            Ok(p) => {
                tracing::info!("Successfully connected to Postgres at {}", &config.url);
                p
            }
            Err(e) => {
                tracing::error!("Failed to connect to Postgres at {}: {}", &config.url, e);
                panic!("Failed to connect to Postgres: {}", e);
            }
        };

        Self { pool }
    }
}
