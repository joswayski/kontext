pub fn load_env() {
    match dotenvy::dotenv().ok() {
        Some(_) => println!("Found .env!"),
        None => println!(".env file not found!"),
    }
}
