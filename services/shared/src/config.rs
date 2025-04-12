use dotenvy::dotenv;
use std::{env, path::PathBuf};

pub fn load_env() -> Option<PathBuf> {
    dotenv().ok()
}

pub fn get_port() -> String {
    env::var("PORT").unwrap_or_else(|_| "4000".to_string())
}
