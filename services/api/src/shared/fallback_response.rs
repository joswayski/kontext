use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, PartialEq)]

pub struct FallbackResponse {
    pub path: String,
    pub method: String,
    pub status: String,
    pub code: u16,
    pub message: String,
    pub documentation: String,
}

impl FallbackResponse {
    pub fn new(
        path: String,
        method: String,
        status: String,
        code: u16,
        message: String,
        documentation: String,
    ) -> Self {
        Self {
            path,
            method,
            status,
            code,
            message,
            documentation,
        }
    }
}
