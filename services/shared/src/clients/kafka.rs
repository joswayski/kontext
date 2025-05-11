use rdkafka::{
    admin::{AdminClient, AdminOptions, NewTopic, TopicReplication},
    client::DefaultClientContext,
    config::{ClientConfig, FromClientConfig},
    consumer::{BaseConsumer, Consumer},
    error::KafkaError,
};

use std::time::Duration;
use tracing;

use crate::config::kafka::{KafkaClusterConfig, KafkaConfig};

struct ClusterClients {
    config: KafkaClusterConfig,
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

        Ok(Self {
            config: config.clone(),
            consumer,
            admin,
        })
    }
}

pub struct KafkaClient {
    clusters: std::collections::HashMap<String, ClusterClients>,
}

impl KafkaClient {
    pub fn new(config: KafkaConfig) -> Result<Self, KafkaError> {
        let mut clusters = std::collections::HashMap::new();

        for (name, cluster_config) in config.clusters {
            let clients = ClusterClients::new(&name, &cluster_config)?;
            clusters.insert(name, clients);
        }

        Ok(Self { clusters })
    }

    fn get_cluster(&self, cluster_name: &str) -> Result<&ClusterClients, KafkaError> {
        self.clusters.get(cluster_name).ok_or_else(|| {
            let err = format!("Cluster {} not found", cluster_name);
            tracing::error!("{}", err);
            KafkaError::ClientCreation(err)
        })
    }

    /// Lists all topics in a cluster
    pub async fn list_topics(&self, cluster_name: &str) -> Result<Vec<String>, KafkaError> {
        let cluster = self.get_cluster(cluster_name)?;
        let metadata = cluster
            .consumer
            .fetch_metadata(None, Duration::from_secs(10))?;

        Ok(metadata
            .topics()
            .iter()
            .map(|topic| topic.name().to_string())
            .collect())
    }

    /// Gets consumer group information for a cluster
    pub async fn list_consumer_groups(
        &self,
        cluster_name: &str,
    ) -> Result<Vec<String>, KafkaError> {
        // For now, just return our known consumer group
        // TODO: Implement proper consumer group listing using the admin client
        Ok(vec![format!("kontext-{}", cluster_name)])
    }

    /// Gets detailed information about a specific topic
    pub async fn describe_topic(
        &self,
        cluster_name: &str,
        topic_name: &str,
    ) -> Result<(), KafkaError> {
        let cluster = self.get_cluster(cluster_name)?;
        let metadata = cluster
            .consumer
            .fetch_metadata(Some(topic_name), Duration::from_secs(10))?;

        if let Some(topic) = metadata.topics().iter().find(|t| t.name() == topic_name) {
            tracing::info!("Topic {} details:", topic_name);
            tracing::info!("  Partitions: {}", topic.partitions().len());
            for partition in topic.partitions() {
                tracing::info!(
                    "    Partition {}: Leader {}",
                    partition.id(),
                    partition.leader()
                );
            }
        }

        Ok(())
    }

    /// Creates a new topic
    pub async fn create_topic(
        &self,
        cluster_name: &str,
        topic_name: &str,
        partitions: i32,
        replication_factor: i32,
    ) -> Result<(), KafkaError> {
        let cluster = self.get_cluster(cluster_name)?;
        let topic = NewTopic::new(
            topic_name,
            partitions,
            TopicReplication::Fixed(replication_factor),
        );
        let results = cluster
            .admin
            .create_topics(&[topic], &AdminOptions::new())
            .await?;

        for result in results {
            match result {
                Ok(_) => tracing::info!("Successfully created topic {}", topic_name),
                Err((name, e)) => tracing::error!("Failed to create topic {}: {}", name, e),
            }
        }

        Ok(())
    }
}
