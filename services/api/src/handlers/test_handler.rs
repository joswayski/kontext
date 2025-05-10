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

pub async fn test_handler(Json(payload): Json<CreateUserInput>) -> Result<impl IntoResponse, AppError> {
    // --- Validation Step ---
    // The `?` operator will convert ValidationErrors into AppError::Validation
    // thanks to the `From<ValidationErrors> for AppError` implementation.
    // Axum will then use the `IntoResponse for AppError` implementation to create the HTTP response.
    payload.validate()?;

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

// This implementation tells Axum how to convert AppError into a response.
impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        match self {
            AppError::Json(rejection) => {
                let status = match rejection {
                    JsonRejection::JsonDataError(_) => StatusCode::BAD_REQUEST,
                    JsonRejection::JsonSyntaxError(_) => StatusCode::BAD_REQUEST,
                    JsonRejection::MissingJsonContentType(_) => StatusCode::UNSUPPORTED_MEDIA_TYPE,
                    // axum::extract::rejection::JsonRejection is non_exhaustive.
                    // This handles other variants, including potential future ones or BytesRejection
                    // from the 'json-with-graceful-shutdown' feature.
                    _ => StatusCode::BAD_REQUEST,
                };
                // The `Display` impl for `JsonRejection` provides a user-friendly error message,
                // e.g., "Failed to deserialize the JSON body into the target type: missing field `email` at line X column Y".
                let body = Json(json!({
                    "message": rejection.to_string()
                    // If you want to add other fields, you can do so here, for example:
                    // "error_type": "JsonProcessingError",
                    // "request_id": "some-uuid" // If you have request tracing
                }));
                (status, body).into_response()
            }
            AppError::Validation(validation_errors) => {
                // Convert validation errors into a structured JSON response
                let errors = validation_errors
                    .field_errors()
                    .into_iter()
                    .map(|(field, field_errors_list)| {
                        let messages: Vec<String> = field_errors_list
                            .iter()
                            .map(|e| {
                                e.message
                                    .as_ref()
                                    .map(|m| m.to_string())
                                    .unwrap_or_else(|| {
                                        // Provide a default message if none is set in the validator
                                        format!("Validation failed for rule: {:?}", e.code)
                                    })
                            })
                            .collect();
                        (field.to_string(), messages) // Ensure field name is a String for HashMap key
                    })
                    .collect::<std::collections::HashMap<String, Vec<String>>>();

                // Return a 400 Bad Request with the validation errors
                (
                    StatusCode::BAD_REQUEST,
                    Json(json!({
                        "message": "Input validation failed",
                        "errors": errors
                    })),
                )
                    .into_response()
            }
        }
    }
}
