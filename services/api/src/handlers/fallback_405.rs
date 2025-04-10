use crate::shared::FallbackResponse;
use axum::{
    extract::Request,
    http::{Method, StatusCode},
    response::IntoResponse,
    Json,
};
pub async fn fallback_405(method: Method, request: Request) -> impl IntoResponse {
    let path = request.uri().path().to_string();
    let response = FallbackResponse::new(
        path,
        method.to_string(),
        "Method Not Allowed".to_string(),
        405,
        "The method you specified is not allowed".to_string(),
        "https://josevalerio.com/kontext".to_string(),
    );

    (StatusCode::METHOD_NOT_ALLOWED, Json(response))
}
