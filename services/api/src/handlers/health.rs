use axum::response::IntoResponse;
use serde::Serialize;

#[derive(Serialize)]
struct HealthResponse {
    status: String,
}

pub async fn health_check() -> impl IntoResponse {
    let response = HealthResponse {
        status: "ok".to_string(),
    };

    axum::Json(response)
}
