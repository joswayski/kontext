use rdkafka::{consumer::BaseConsumer, ClientConfig};

struct KafkaClients {
    pub consumer: BaseConsumer,
}

impl KafkaClients {
    pub fn new_consumer(brokers: String) -> BaseConsumer {
        ClientConfig::new()
            .set("bootstrap.servers", brokers)
            .create()
            .expect("Consumer creation failed")
    }
}
