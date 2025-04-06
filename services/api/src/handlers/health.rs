use axum::response::IntoResponse;
use serde::Serialize;

#[derive(Serialize)]
struct HealthResponse {
    status: String,
}

pub async fn health_check() -> impl IntoResponse {
    let response = HealthResponse {
        status: "Saul Goodman".to_string(),
    };

    axum::Json(response)
}
