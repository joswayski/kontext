use serde::{Deserialize, Serialize};
use std::{collections::HashMap, env};

// A single cluster config
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct KafkaClusterConfig {
    pub brokers: String,
    pub metrics_url: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct KafkaConfig {
    // key is the cluster name
    pub clusters: HashMap<String, KafkaClusterConfig>,
}

impl KafkaConfig {
    pub fn new() -> Self {
        let mut clusters = HashMap::new();

        // Get all environment variables that start with KAFKA_
        let env_vars: HashMap<String, String> = env::vars()
            .filter(|(key, _)| key.starts_with("KAFKA_"))
            .collect();

        // Find all unique cluster names (KAFKA_{CLUSTERNAME}_BROKERS)
        let cluster_names: Vec<String> = env_vars
            .keys()
            .filter_map(|key| {
                key.strip_prefix("KAFKA_")
                    .and_then(|s| s.split('_').next())
                    .map(String::from)
            })
            .collect::<std::collections::HashSet<_>>()
            .into_iter()
            .collect();

        for cluster_name in cluster_names {
            let brokers = env::var(format!("KAFKA_{}_BROKERS", cluster_name))
                .unwrap_or_else(|_| "localhost:9092".to_string());

            let metrics_url = env::var(format!("KAFKA_{}_METRICS_URL", cluster_name))
                .unwrap_or_else(|_| format!("http://localhost:8080/metrics/{}", cluster_name));

            clusters.insert(
                cluster_name.clone(),
                KafkaClusterConfig {
                    brokers,
                    metrics_url,
                },
            );
        }

        Self { clusters }
    }

    pub fn get_cluster(&self, name: &str) -> Option<&KafkaClusterConfig> {
        self.clusters.get(name)
    }
}
