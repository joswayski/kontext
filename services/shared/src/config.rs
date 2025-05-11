use dotenvy::dotenv;
use kafka::KafkaConfig;
use std::env;

pub mod kafka;
pub mod postgres;
pub mod tracing;

use postgres::PostgresConfig;

#[derive(Debug, Clone)]
pub struct Config {
    pub port: String,
    pub kafka: KafkaConfig,
    pub tracing: tracing::TracingConfig,
    pub postgres: PostgresConfig,
}

impl Config {
    /**
     * Starts up a config with a default port, kafka clients, and tracing config.
     * TODO allow overriding settings in the future
     * Will panic if required connections (Kafka clusters, Postgres) cannot be established.
     */
    pub fn init() -> Self {
        dotenv().ok();

        let config = Self {
            port: env::var("PORT").unwrap_or_else(|_| "4000".to_string()),
            kafka: KafkaConfig::new(),
            tracing: tracing::TracingConfig::default(),
            postgres: PostgresConfig::new(),
        };

        // Only initialize tracing if it hasn't been initialized yet
        static TRACING_INITIALIZED: std::sync::Once = std::sync::Once::new();
        TRACING_INITIALIZED.call_once(|| {
            config.tracing.init();
        });

        config
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use serial_test::serial;
    use std::env;

    #[test]
    #[serial]
    fn test_config_defaults() {
        // Clear any existing env vars
        env::remove_var("PORT");
        env::remove_var("KAFKA_CLUSTER1_BROKER_URL");
        env::remove_var("KAFKA_CLUSTER2_BROKER_URL");

        let config = Config::init();

        // Default port
        assert_eq!(config.port, "4000");

        // Tracing
        assert_eq!(config.tracing.level, tracing::Level::INFO);
        assert!(config.tracing.with_ansi);
        assert!(config.tracing.with_level);
        assert_eq!(config.kafka.clusters.capacity(), 0); // no clusters by default
    }

    #[test]
    #[serial]
    fn test_config_from_env() {
        // First clear any existing env vars to ensure clean state
        env::remove_var("PORT");
        env::remove_var("KAFKA_CLUSTER1_BROKER_URL");
        env::remove_var("KAFKA_CLUSTER2_BROKER_URL");

        // Then set test values
        env::set_var("PORT", "69420");
        env::set_var("KAFKA_CLUSTER1_BROKER_URL", "kafka1:9092");
        env::set_var("KAFKA_CLUSTER2_BROKER_URL", "kafka2:9092,kafka3:9092");

        let config = Config::init();

        // Test port
        assert_eq!(config.port, "69420");

        // Test Kafka clusters
        let cluster1 = config.kafka.get_cluster("CLUSTER1").unwrap();
        assert_eq!(cluster1.brokers, "kafka1:9092");

        let cluster2 = config.kafka.get_cluster("CLUSTER2").unwrap();
        assert_eq!(cluster2.brokers, "kafka2:9092,kafka3:9092");

        // Clean up
        env::remove_var("PORT");
        env::remove_var("KAFKA_CLUSTER1_BROKER_URL");
        env::remove_var("KAFKA_CLUSTER2_BROKER_URL");
    }
}
