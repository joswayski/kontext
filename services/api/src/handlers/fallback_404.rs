use crate::shared::fallback_response::FallbackResponse;
use axum::{
    extract::Request,
    http::{Method, StatusCode},
    response::IntoResponse,
    Json,
};

pub async fn fallback_404(method: Method, request: Request) -> impl IntoResponse {
    let path = request.uri().path().to_string();
    let response = FallbackResponse::new(
        path,
        method.to_string(),
        "Not Found".to_string(),
        404,
        "The route you specified is not found".to_string(),
        "https://josevalerio.com/kontext".to_string(),
    );

    (StatusCode::NOT_FOUND, Json(response))
}
