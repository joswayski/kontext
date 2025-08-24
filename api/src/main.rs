use shared::load_env;
use std::env;
fn main() {
    load_env();
    println!("Hello from API!");

    for (key, value) in env::vars() {
        println!("Key: {} - Value: {}", key, value)
    }
}
