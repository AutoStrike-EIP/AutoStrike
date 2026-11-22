#[derive(Debug)]
pub struct AddressParams {
    pub host: String,
    pub port: u16,
    pub user: String,
    pub password: String,
    pub vhost: String,
}

impl AddressParams {
    pub fn to_string(&self) -> String {
        format!("amqp://{}:{}@{}:{}/{}", self.user, self.password, self.host, self.port, self.vhost)
    }
}

impl Default for AddressParams {
    fn default() -> Self {
        Self {
            host: "127.0.0.1".to_string(),
            port: 5672,
            user: "guest".to_string(),
            password: "guest".to_string(),
            vhost: "%2f".to_string(),
        }
    }
}
