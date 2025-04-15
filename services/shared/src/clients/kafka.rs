use rdkafka::{consumer::BaseConsumer, ClientConfig};

struct KafkaClients {
    pub consumer: BaseConsumer,
}

impl KafkaClients {
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
}
