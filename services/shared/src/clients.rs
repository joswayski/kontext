use qdrant_client::prelude::*;
use rdkafka::{consumer::BaseConsumer, ClientConfig};
use sqlx::{mysql::MySqlPoolOptions, MySql, Pool};

pub struct AllClients {
    pub kafka: BaseConsumer,
    pub mysql: Pool<MySql>,
    pub qdrant: QdrantClient,
}

impl AllClients {
    pub fn create_consumer(brokers: String) -> Option<BaseConsumer> {
        match ClientConfig::new()
            .set("bootstrap.servers", &brokers)
            .create()
        {
            Ok(consumer) => {
                tracing::info!("Consumer created for {}", &brokers);
                Some(consumer)
            }
            Err(e) => {
                tracing::error!("Failed to create consumer on {} - {}", brokers, e);
                None
            }
        }
    }

    pub async fn create_mysql_pool(database_url: &str) -> Option<Pool<MySql>> {
        match MySqlPoolOptions::new()
            .max_connections(5)
            .connect(database_url)
            .await
        {
            Ok(pool) => {
                tracing::info!("MySQL pool created successfully");
                Some(pool)
            }
            Err(e) => {
                tracing::error!("Failed to create MySQL pool: {}", e);
                None
            }
        }
    }

    pub async fn create_qdrant_client(qdrant_url: &str) -> Option<QdrantClient> {
        match QdrantClient::from_url(qdrant_url).build() {
            Ok(client) => {
                tracing::info!("Qdrant client created successfully");
                Some(client)
            }
            Err(e) => {
                tracing::error!("Failed to create Qdrant client: {}", e);
                None
            }
        }
    }

    pub async fn init(
        kafka_brokers: String,
        mysql_url: String,
        qdrant_url: String,
    ) -> Option<Self> {
        let kafka = Self::create_consumer(kafka_brokers)?;
        let mysql = Self::create_mysql_pool(&mysql_url).await?;
        let qdrant = Self::create_qdrant_client(&qdrant_url).await?;

        Self {
            kafka,
            mysql,
            qdrant,
        }
    }
}
