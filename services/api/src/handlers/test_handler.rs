use axum::{extract::rejection::JsonRejection, http::StatusCode, response::IntoResponse, Json};
use serde::Deserialize;
use serde_json::json;
use validator::{Validate, ValidationErrors};

#[derive(Deserialize, Validate)]
pub struct CreateUserInput {
    #[validate(email(message = "Please provide a valid email address."))]
    pub email: String,

    #[validate(length(min = 8, message = "Password must be at least 8 characters long."))]
    pub password: String,

    #[validate(range(min = 18, max = 120, message = "Age must be between 18 and 120."))]
    pub age: u32,

    #[validate(url(message = "Please provide a valid URL for your website."))]
    pub website: Option<String>, // Can validate Option fields

    #[validate(must_match(other = "password", message = "Password confirmation does not match."))]
    pub password_confirmation: String,
}

pub async fn test_handler(
    payload: Result<Json<CreateUserInput>, JsonRejection>,
) -> impl IntoResponse {
    let payload = match payload {
        Ok(payload) => payload.0,
        Err(err) => {
            return (
                StatusCode::BAD_REQUEST,
                Json(json!({
                    "error": "Invalid request body",
                    "message": "Failed to parse request body. Please ensure all required fields are provided.",
                    "details": err.to_string(),
                    "documentation": "https://josevalerio.com/kontext/api"
                })),
            );
        }
    };

    if let Err(errors) = payload.validate() {
        return (
            StatusCode::BAD_REQUEST,
            Json(json!({
                "error": "Validation failed",
                "message": "Request validation failed. Please check the details.",
                "details": errors,
                "documentation": "https://your-api-docs.com/schemas/create-user"
            })),
        );
    }

    println!(
        "User creation request validated successfully for email: {}",
        payload.email
    );

    (
        StatusCode::CREATED,
        Json(json!({"message": "User created"})),
    )
}
