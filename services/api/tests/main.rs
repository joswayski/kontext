use axum::{
    body::Body,
    http::{Request, StatusCode},
};
use tower::ServiceExt; // Needed for .oneshot()
async fn get_response(uri: &str) -> axum::response::Response {
    // Call the function from the library crate 'api'
    let app = create_routes();
    let request = Request::builder().uri(uri).body(Body::empty()).unwrap();

    app.oneshot(request).await.unwrap()
}

#[tokio::test]
async fn test_root_route() {
    let response = get_response("/").await;
    assert_eq!(response.status(), StatusCode::OK);
    let body = axum::body::to_bytes(response.into_body(), usize::MAX)
        .await
        .unwrap();
    assert_eq!(
        &body[..],
        b"Welcome to Kontext API :) - Check the docs for more information!"
    );
}

#[tokio::test]
async fn test_health_check_route() {
    let response = get_response("/health").await;
    assert_eq!(response.status(), StatusCode::OK);
    // TODO: Assert the body of the response
}

#[tokio::test]
async fn test_api_health_check_route() {
    let response = get_response("/api/health").await;
    assert_eq!(response.status(), StatusCode::OK);
    // TODO: Assert the body of the response
}

#[tokio::test]
async fn test_not_found_route() {
    let uri = "/a/route/that/does/not/exist";
    let response = get_response(uri).await;

    assert_eq!(response.status(), StatusCode::NOT_FOUND);

    // Extract bytes from the body
    let body_bytes = axum::body::to_bytes(response.into_body(), usize::MAX)
        .await
        .unwrap();

    // Deserialize the body into your FallbackResponse struct
    let fallback_resp: FallbackResponse = serde_json::from_slice(&body_bytes)
        .expect("Failed to deserialize response body into FallbackResponse");

    // Assert the content of the response
    assert_eq!(fallback_resp.path, uri);
    assert_eq!(fallback_resp.method, "GET");
    assert_eq!(fallback_resp.status, "Not Found");
    assert_eq!(fallback_resp.code, 404);
    assert_eq!(
        fallback_resp.message,
        "The route you specified is not found"
    );
    assert_eq!(
        fallback_resp.documentation,
        "https://josevalerio.com/kontext"
    );
}

// Add tests for 405 Method Not Allowed if needed
#[tokio::test]
async fn test_method_not_allowed() {
    let app = create_routes();
    // Try POSTing to a GET-only route like /health
    let request = Request::builder()
        .uri("/health")
        .method("POST")
        .body(Body::empty())
        .unwrap();

    let response = app.oneshot(request).await.unwrap();
    assert_eq!(response.status(), StatusCode::METHOD_NOT_ALLOWED);
    // Assert the body of the 405 if you have custom content
}
