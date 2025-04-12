use axum::response::IntoResponse;
use serde::Serialize;

#[derive(Serialize)]
struct RootResponse {
    message: String,
}

pub async fn root() -> impl IntoResponse {
    let response = RootResponse {
        message: "Welcome to Kontext API :) - Check the docs for more information!".to_string(),
    };

    axum::Json(response)
}
