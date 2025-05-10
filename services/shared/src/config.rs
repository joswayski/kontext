use dotenvy::dotenv;
use kafka::KafkaConfig;
use std::env;

pub mod kafka;
pub mod tracing;

#[derive(Debug, Clone)]
pub struct Config {
    pub port: String,
    pub kafka: KafkaConfig,
    pub tracing: tracing::TracingConfig,
}

impl Config {
    /**
     * Starts up a config with a default port, kafka clients, and tracing config.
     * TODO allow overriding settings in the future
     */
    pub fn init() -> Self {
        dotenv().ok();

        let config = Self {
            port: env::var("PORT").unwrap_or_else(|_| "4000".to_string()),
            kafka: KafkaConfig::new(),
            tracing: tracing::TracingConfig::default(),
        };

        config.tracing.init();

        config
    }
}

#[test]
fn test_config_defaults() {
    // Clear any existing env vars
    env::remove_var("PORT");
    env::remove_var("KAFKA_CLUSTER1_BROKERS");
    env::remove_var("KAFKA_CLUSTER1_METRICS_URL");

    let config = Config::init();

    // Default port
    assert_eq!(config.port, "4000");

    // Tracing
    assert_eq!(config.tracing.level, tracing::Level::DEBUG);
    assert!(config.tracing.with_ansi);
    assert!(config.tracing.with_level);
    assert_eq!(config.kafka.clusters.capacity(), 0); // no clusters by default
}

#[test]
fn test_config_from_env() {
    // Set test values
    env::set_var("PORT", "8080");
    env::set_var("KAFKA_CLUSTER1_BROKERS", "kafka1:9092");
    env::set_var("KAFKA_CLUSTER1_METRICS_URL", "http://cluster1:8080/metrics");
    env::set_var("KAFKA_CLUSTER2_BROKERS", "kafka2:9092,kafka3:9092");
    env::set_var("KAFKA_CLUSTER2_METRICS_URL", "http://cluster2:8080/metrics");

    let config = Config::init();

    // Test port
    assert_eq!(config.port, "8080");

    // Test Kafka clusters
    let cluster1 = config.kafka.get_cluster("CLUSTER1").unwrap();
    assert_eq!(cluster1.brokers, "kafka1:9092");
    // TODO
    // assert_eq!(cluster1.metrics_url, "http://cluster1:8080/metrics");

    let cluster2 = config.kafka.get_cluster("CLUSTER2").unwrap();
    assert_eq!(cluster2.brokers, "kafka2:9092,kafka3:9092");
    // TODO
    // assert_eq!(cluster2.metrics_url, "http://cluster2:8080/metrics");

    // Clean up
    env::remove_var("PORT");
    env::remove_var("KAFKA_CLUSTER1_BROKERS");
    env::remove_var("KAFKA_CLUSTER1_METRICS_URL");
    env::remove_var("KAFKA_CLUSTER2_BROKERS");
    env::remove_var("KAFKA_CLUSTER2_METRICS_URL");
}
