use axum::{extract::rejection::JsonRejection, http::StatusCode, response::IntoResponse, Json};
use serde::{Deserialize, Serialize};
use serde_json::json;
use std::collections::HashMap;

#[derive(Deserialize, Serialize)]
pub struct CreateUserInput {
    pub email: String,
    pub password: String,
    pub age: u32,
    pub website: Option<String>,
    pub password_confirmation: String,
}

#[derive(Serialize)]
pub struct FieldError {
    pub code: &'static str,
    pub message: String,
}

impl CreateUserInput {
    pub fn validate(&self) -> Result<(), HashMap<&'static str, Vec<FieldError>>> {
        let mut errors: HashMap<&'static str, Vec<FieldError>> = HashMap::new();

        if !self.email.contains('@') {
            errors.entry("email").or_default().push(FieldError {
                code: "email",
                message: "Please provide a valid email address.".to_string(),
            });
        }

        if self.password.len() < 8 {
            errors.entry("password").or_default().push(FieldError {
                code: "length",
                message: "Password must be at least 8 characters long.".to_string(),
            });
        }

        if self.password != self.password_confirmation {
            errors
                .entry("password_confirmation")
                .or_default()
                .push(FieldError {
                    code: "must_match",
                    message: "Password confirmation does not match.".to_string(),
                });
        }

        if !(18..=120).contains(&self.age) {
            errors.entry("age").or_default().push(FieldError {
                code: "range",
                message: "Age must be between 18 and 120.".to_string(),
            });
        }

        if let Some(website_url) = &self.website {
            if !website_url.starts_with("http://") && !website_url.starts_with("https://") {
                errors.entry("website").or_default().push(FieldError {
                    code: "url",
                    message: "Please provide a valid URL for your website.".to_string(),
                });
            }
        }

        if errors.is_empty() {
            Ok(())
        } else {
            Err(errors)
        }
    }
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
            })),
        );
    }

    (
        StatusCode::CREATED,
        Json(json!({"message": "User created", "user": &payload})),
    )
}
