use rdkafka::{
    admin::AdminClient,
    client::DefaultClientContext,
    config::{ClientConfig, FromClientConfig},
    consumer::BaseConsumer,
    error::KafkaError,
};

use tracing;

use crate::config::kafka::{KafkaClusterConfig, KafkaConfig};

struct ClusterClients {
    consumer: BaseConsumer,
    admin: AdminClient<DefaultClientContext>,
}

impl ClusterClients {
    fn new(cluster_name: &str, config: &KafkaClusterConfig) -> Result<Self, KafkaError> {
        // Create consumer
        let mut consumer_config = ClientConfig::new();
        consumer_config
            .set("bootstrap.servers", &config.brokers)
            .set("group.id", format!("kontext-{}", cluster_name))
            .set("enable.auto.commit", "false")
            .set("auto.offset.reset", "earliest");

        let consumer: BaseConsumer = consumer_config.create().map_err(|e| {
            tracing::error!(
                "Failed to create Kafka consumer for cluster {}: {}",
                cluster_name,
                e
            );
            e
        })?;

        // Create admin client
        let mut admin_config = ClientConfig::new();
        admin_config.set("bootstrap.servers", &config.brokers);

        let admin = AdminClient::from_config(&admin_config).map_err(|e| {
            tracing::error!(
                "Failed to create Kafka admin client for cluster {}: {}",
                cluster_name,
                e
            );
            e
        })?;

        Ok(Self { consumer, admin })
    }
}

pub struct KafkaClient {
    clusters: std::collections::HashMap<String, ClusterClients>,
}

impl KafkaClient {
    pub fn new(config: KafkaConfig) -> Self {
        let mut cluster_clients_map = std::collections::HashMap::new();

        for (cluster_name, cluster_config) in config.clusters {
            match ClusterClients::new(&cluster_name, &cluster_config) {
                Ok(clients) => {
                    tracing::info!(
                        "Successfully created Kafka clients for cluster {}",
                        cluster_name
                    );
                    cluster_clients_map.insert(cluster_name.clone(), clients);
                }
                Err(e) => {
                    // Panic here if a cluster client fails to initialize
                    panic!(
                        "Failed to create Kafka clients for cluster {}: {}. Aborting KafkaClient creation.",
                        cluster_name,
                        e
                    );
                }
            }
        }

        Self {
            clusters: cluster_clients_map,
        }
    }

    pub fn get_admin_client(
        &self,
        cluster_name: &str,
    ) -> Option<&AdminClient<DefaultClientContext>> {
        self.clusters.get(cluster_name).map(|cc| &cc.admin)
    }

    pub fn get_consumer(&self, cluster_name: &str) -> Option<&BaseConsumer> {
        self.clusters.get(cluster_name).map(|cc| &cc.consumer)
    }
}
