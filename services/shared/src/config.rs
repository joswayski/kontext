use dotenvy::dotenv;
use std::env;

pub fn load_env() {
    dotenv().ok();
}

pub fn get_port() -> String {
    env::var("PORT").unwrap_or_else(|_| "4000".to_string())
}
