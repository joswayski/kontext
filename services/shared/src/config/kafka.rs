use dotenvy::dotenv;
use serde::{Deserialize, Serialize};
use std::{collections::HashMap, env};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct KafkaClusterConfig {
    pub brokers: String,
    // pub metrics_url: String, // TODO
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct KafkaConfig {
    // key is the cluster name
    pub clusters: HashMap<String, KafkaClusterConfig>,
}

const PREFIX: &str = "KAFKA_";
const BROKER_SUFFIX: &str = "_BROKER_URL";
// const METRICS_SUFFIX: &str = "_METRICS_URL";
const ALL_ENDINGS: [&str; 1] = [
    BROKER_SUFFIX,
    // METRICS_SUFFIX
]; // TODO perhaps better as enum

impl KafkaConfig {
    pub fn new() -> Self {
        dotenv().ok();

        let mut clusters: HashMap<String, KafkaClusterConfig> = HashMap::new();

        for (key, value) in env::vars() {
            let rest = match key.strip_prefix(PREFIX) {
                Some(rest) => rest,
                None => continue, // Skip keys that dont have this
            };

            let (cluster_name, suffix) = match ALL_ENDINGS
                .iter()
                .find_map(|ending| rest.strip_suffix(ending).map(|name| (name, *ending)))
            {
                Some((name, suffix)) => (name, suffix),
                None => continue, // Skip keys that don't end in one of our suffixes
            };

            let cluster_config = clusters.entry(cluster_name.to_string()).or_default();

            match suffix {
                BROKER_SUFFIX => {
                    cluster_config.brokers = value;
                }
                // METRICS_SUFFIX => {
                //     cluster_config.metrics_url = value;
                // }
                _ => {}
            }
        }

        Self { clusters }
    }

    pub fn get_cluster(&self, name: &str) -> Option<&KafkaClusterConfig> {
        self.clusters.get(name)
    }
}

impl Default for KafkaClusterConfig {
    fn default() -> Self {
        Self {
            brokers: String::from("BROKER URL NOT SET"),
            // metrics_url: String::from("METRICS URL NOT SET"),
        }
    }
}
