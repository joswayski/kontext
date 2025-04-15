use rdkafka::{consumer::BaseConsumer, ClientConfig};

use crate::config::kafka::KafkaConfig;

struct KafkaClients {
    pub consumer: BaseConsumer,
}

impl KafkaClients {
    pub fn new_consumer(brokers: String, config: KafkaConfig) -> Option<BaseConsumer> {
        match ClientConfig::new()
            .set("bootstrap.servers", brokers)
            .create()
        {
            Ok(consumer) => consumer,
            Err(e) => {
                tracing::error!("Failed to create consumer on {}", brokers);
                None
            }
        }
    }
}
