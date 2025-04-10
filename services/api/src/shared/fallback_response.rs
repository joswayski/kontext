use serde::Serialize;

#[derive(Serialize)]
pub struct FallbackResponse {
    path: String,
    method: String,
    status: String,
    code: u16,
    message: String,
    documentation: String,
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
