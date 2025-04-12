use axum::response::IntoResponse;
use serde::Serialize;

#[derive(Serialize)]
struct HealthResponse {
    message: String,
}

pub async fn health_check() -> impl IntoResponse {
    let response = HealthResponse {
        message: "Saul Goodman".to_string(),
    };

    axum::Json(response)
}
