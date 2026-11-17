pub struct AddressParams {
    pub host: String,
    pub port: u16,
    pub user: String,
    pub password: String,
    pub vhost: String,
}

impl AddressParams {
    pub fn default() -> Self {
        AddressParams {
            host: "127.0.0.1".to_string(),
            port: 5672,
            user: "guest".to_string(),
            password: "guest".to_string(),
            vhost: "%2f".to_string(),
        }
    }
}
