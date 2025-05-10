use axum::{
    extract::rejection::JsonRejection,
    http::StatusCode,
    response::{IntoResponse, Response},
    Json,
};
use serde::Deserialize;
use serde_json::json;

use validator::Validate;

use validator::ValidationErrors; // Make sure this is imported

#[derive(Debug)] // Add Debug for easier troubleshooting
pub enum AppError {
    // Variant to wrap JSON parsing/extraction errors
    Json(JsonRejection),
    // Variant to wrap your validation errors
    Validation(ValidationErrors),
    // You could add other error types here later
}

#[derive(Deserialize, Validate)] // Add Validate here
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

pub async fn test_handler(Json(payload): Json<CreateUserInput>) -> impl IntoResponse {
    // --- Validation Step ---
    if let Err(validation_errors) = payload.validate() {
        // Validation failed, return a user-friendly error response
        println!("Validation failed: {:?}", validation_errors); // Log details server-side

        // Convert validation errors into a structured JSON response
        // The structure here is up to you, but field-level errors are common.
        let errors = validation_errors
            .field_errors()
            .into_iter()
            .map(|(field, errors)| {
                let messages: Vec<String> = errors
                    .iter()
                    .map(|e| {
                        // Use the custom message if provided, otherwise format a default
                        e.message
                            .as_ref()
                            .map(|m| m.to_string())
                            .unwrap_or_else(|| format!("Validation failed for rule: {:?}", e.code))
                    })
                    .collect();
                (field, messages)
            })
            .collect::<std::collections::HashMap<_, _>>();

        // Return a 400 Bad Request with the errors
        return Err((
            StatusCode::BAD_REQUEST,
            Json(json!({
                "message": "Input validation failed",
                "errors": errors
            })),
        ));
    }

    // --- Validation Passed ---
    // Proceed with your application logic (e.g., save the user)
    println!(
        "User creation request validated successfully for email: {}",
        payload.email
    );

    // Return a success response (e.g., 201 Created)
    Ok((
        StatusCode::CREATED,
        Json(json!({"message": "User created"})),
    ))
}

impl From<JsonRejection> for AppError {
    fn from(rejection: JsonRejection) -> Self {
        AppError::Json(rejection)
    }
}

// Convert ValidationErrors into AppError::Validation
impl From<ValidationErrors> for AppError {
    fn from(errors: ValidationErrors) -> Self {
        AppError::Validation(errors)
    }
}

// This implementation tells Axum how to convert a JSON rejection into a response.
impl IntoResponse for JsonRejection {
    fn into_response(self) -> Response {
        let (status, message) = match self {
            JsonRejection::JsonDataError(e) => {
                // Error deserializing - e.g., wrong types
                (StatusCode::BAD_REQUEST, format!("Invalid JSON data: {}", e))
            }
            JsonRejection::JsonSyntaxError(e) => {
                // Error in JSON syntax - e.g., trailing comma, missing quotes
                (
                    StatusCode::BAD_REQUEST,
                    format!("Invalid JSON syntax: {}", e),
                )
            }
            JsonRejection::MissingJsonContentType(_) => {
                // Missing Content-Type header
                (
                    StatusCode::UNSUPPORTED_MEDIA_TYPE,
                    "Missing 'Content-Type: application/json' header".to_string(),
                )
            }
            // Catch-all for other potential JSON rejection reasons
            _ => (
                StatusCode::BAD_REQUEST,
                format!("Unknown JSON error: {}", self),
            ),
        };

        // Simple JSON response
        let body = Json(json!({
            "message": "Failed to process request body", // General message
            "error": message // Specific reason
        }));

        (status, body).into_response()
    }
}
