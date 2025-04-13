use dotenvy::dotenv;
use kafka::KafkaConfig;
use std::env;

mod kafka;

#[derive(Debug, Clone)]
pub struct Config {
    pub port: String,
    pub kafka: KafkaConfig,
}

impl Config {
    pub fn new() -> Self {
        // Load environment variables from .env file if it exists
        let _ = dotenv();

        Self {
            port: env::var("PORT").unwrap_or_else(|_| "4000".to_string()),
            kafka: KafkaConfig::new(),
        }
    }
}

#[test]
fn test_config_defaults() {
    // Clear any existing env vars
    env::remove_var("PORT");
    env::remove_var("KAFKA_CLUSTER1_BROKERS");
    env::remove_var("KAFKA_CLUSTER1_METRICS_URL");

    let config = Config::new();

    // Test default port
    assert_eq!(config.port, "4000");

    // Test default Kafka cluster
    let default_cluster = config.kafka.get_cluster("DEFAULT").unwrap();
    assert_eq!(default_cluster.brokers, "localhost:9092");
    assert_eq!(
        default_cluster.metrics_url,
        "http://localhost:8080/metrics/DEFAULT"
    );
}

#[test]
fn test_config_from_env() {
    // Set test values
    env::set_var("PORT", "8080");
    env::set_var("KAFKA_CLUSTER1_BROKERS", "kafka1:9092");
    env::set_var("KAFKA_CLUSTER1_METRICS_URL", "http://cluster1:8080/metrics");
    env::set_var("KAFKA_CLUSTER2_BROKERS", "kafka2:9092,kafka3:9092");
    env::set_var("KAFKA_CLUSTER2_METRICS_URL", "http://cluster2:8080/metrics");

    let config = Config::new();

    // Test port
    assert_eq!(config.port, "8080");

    // Test Kafka clusters
    let cluster1 = config.kafka.get_cluster("CLUSTER1").unwrap();
    assert_eq!(cluster1.brokers, "kafka1:9092");
    assert_eq!(cluster1.metrics_url, "http://cluster1:8080/metrics");

    let cluster2 = config.kafka.get_cluster("CLUSTER2").unwrap();
    assert_eq!(cluster2.brokers, "kafka2:9092,kafka3:9092");
    assert_eq!(cluster2.metrics_url, "http://cluster2:8080/metrics");

    // Clean up
    env::remove_var("PORT");
    env::remove_var("KAFKA_CLUSTER1_BROKERS");
    env::remove_var("KAFKA_CLUSTER1_METRICS_URL");
    env::remove_var("KAFKA_CLUSTER2_BROKERS");
    env::remove_var("KAFKA_CLUSTER2_METRICS_URL");
}
