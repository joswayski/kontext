use dotenvy::dotenv;
use std::{env, path::PathBuf};

pub fn load_env() -> Option<PathBuf> {
    dotenv().ok()
}

pub fn get_port() -> String {
    env::var("PORT").unwrap_or_else(|_| "4000".to_string())
}

#[test]
fn test_get_port_default() {
    // Clear any existing PORT environment variable
    env::remove_var("PORT");

    // Test default value
    assert_eq!(get_port(), "4000");
}

#[test]
fn test_get_port_from_env() {
    // Set a test PORT value
    env::set_var("PORT", "8080");

    // Test that it uses the environment variable
    assert_eq!(get_port(), "8080");

    env::remove_var("PORT");
}
