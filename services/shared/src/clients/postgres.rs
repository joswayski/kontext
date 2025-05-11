use sqlx::postgres::PgPoolOptions;
use tracing;

use crate::config::postgres::PostgresConfig;

pub struct PostgresClient {
    config: PostgresConfig,
    pool: Option<sqlx::PgPool>,
}

impl PostgresClient {
    pub fn new(config: PostgresConfig) -> Self {
        Self { config, pool: None }
    }

    /// Creates or returns an existing connection pool
    pub async fn get_pool(&mut self) -> Result<&sqlx::PgPool, sqlx::Error> {
        if self.pool.is_none() {
            let pool = PgPoolOptions::new()
                .max_connections(self.config.max_connections)
                .connect(&self.config.url)
                .await
                .map_err(|e| {
                    tracing::error!("Failed to connect to Postgres: {}", e);
                    e
                })?;

            tracing::info!("Successfully connected to Postgres");
            self.pool = Some(pool);
        }

        Ok(self.pool.as_ref().unwrap())
    }
}
