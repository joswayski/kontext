use rdkafka::admin::AdminClient;
use rdkafka::metadata::Metadata;
use rdkafka::ClientConfig;
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;

pub struct KafkaClusterClient {
    pub name: String,
    pub brokers: Vec<String>,
    pub admin_client: AdminClient<rdkafka::client::DefaultClientContext>,
    pub metrics_url: String,
}

pub struct KafkaClusterState {
    pub client: KafkaClusterClient,
    pub metadata: RwLock<Option<Metadata>>,
}

pub struct AppState {
    pub clusters: HashMap<String, Arc<KafkaClusterState>>,
}

impl AppState {
    pub fn new(cluster_configs: Vec<ClusterConfig>) -> Self {
        let clusters = cluster_configs
            .into_iter()
            .map(|cluster| {
                let admin_client: AdminClient<_> = ClientConfig::new()
                    .set("bootstrap.servers", cluster.brokers.join(","))
                    .create()
                    .expect("Kafka admin client creation failed");

                let client = KafkaClusterClient {
                    name: cluster.name.clone(),
                    brokers: cluster.brokers,
                    admin_client,
                    metrics_url: cluster.metrics_url,
                };

                (
                    cluster.name,
                    Arc::new(KafkaClusterState {
                        client,
                        metadata: RwLock::new(None),
                    }),
                )
            })
            .collect();

        Self { clusters }
    }
}
