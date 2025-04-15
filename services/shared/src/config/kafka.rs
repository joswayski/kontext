use dotenvy::dotenv;
use serde::{Deserialize, Serialize};
use std::{collections::HashMap, env};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct KafkaClusterConfig {
    pub brokers: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct KafkaConfig {
    // key is the cluster name
    pub clusters: HashMap<String, KafkaClusterConfig>,
}

const PREFIX: &str = "KAFKA_";
const BROKER_SUFFIX: &str = "_BROKER_URL";
const ALL_ENDINGS: [&str; 1] = [BROKER_SUFFIX];

impl KafkaConfig {
    pub fn new() -> Self {
        dotenv().ok();

        let mut all_clusters: HashMap<String, KafkaClusterConfig> = HashMap::new();

        for (key, value) in env::vars() {
            if !key.starts_with(PREFIX) {
                continue;
            }

            // Skip if the key doesn't end with any of our valid suffixes
            if !ALL_ENDINGS.iter().any(|suffix| key.ends_with(suffix)) {
                continue;
            }

            let key_without_prefix_opt = key.strip_prefix(PREFIX);

            match key_without_prefix_opt {
                Some(rest) => {
                    let cluster_name = ALL_ENDINGS
                        .iter()
                        .find_map(|ending| rest.strip_suffix(ending));

                    match cluster_name {
                        Some(name) => {
                            let cluster_config = KafkaClusterConfig { brokers: value };

                            all_clusters
                                .entry(name.to_string())
                                .or_insert(cluster_config);
                        }
                        None => continue,
                    }
                }
                None => continue,
            }
        }

        Self {
            clusters: all_clusters,
        }
    }

    pub fn get_cluster(&self, name: &str) -> Option<&KafkaClusterConfig> {
        self.clusters.get(name)
    }
}
