use tracing_subscriber::fmt;

pub use tracing::Level;

#[derive(Debug, Clone)]
pub struct TracingConfig {
    pub level: Level,
    pub with_ansi: bool,
    pub with_level: bool,
}

impl Default for TracingConfig {
    fn default() -> Self {
        Self {
            level: Level::INFO,
            with_ansi: true,
            with_level: true,
        }
    }
}

impl TracingConfig {
    pub fn new(level: Level, with_ansi: bool, with_level: bool) -> Self {
        Self {
            level,
            with_ansi,
            with_level,
        }
    }

    pub fn init(&self) {
        fmt()
            .with_max_level(self.level)
            .with_level(self.with_level)
            .with_ansi(self.with_ansi)
            .init();
    }
}
