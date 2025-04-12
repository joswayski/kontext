use axum::response::IntoResponse;
use serde::Serialize;

#[derive(Serialize)]
struct RootResponse {
    status: String,
}

pub async fn root() -> impl IntoResponse {
    let response = RootResponse {
        status: "Welcome to Kontext API :) - Check the docs for more information!".to_string(),
    };

    axum::Json(response)
}
