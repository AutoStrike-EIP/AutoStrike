pub mod address_params;
pub mod route_params;

use lapin::{BasicProperties, Channel, ChannelState, Connection, ConnectionProperties, ConnectionState, ExchangeKind};
use lapin::options::{BasicPublishOptions, ExchangeDeclareOptions, QueueBindOptions, QueueDeclareOptions};
use lapin::types::{FieldTable, ShortString};
pub use address_params::{AddressParams};

pub struct RabbitMQ {
    uri: String,
    connection: Option<Connection>,
    channel: Option<Channel>,
    exchange_name: ShortString,
}

impl RabbitMQ {
    /// Create a new instance of RabbitMQ with a specific queue and a specific address.
    ///
    /// # Arguments
    ///
    /// * `address_params`: The uri parameters of the RabbitMQ connection.
    ///
    /// # Returns
    /// An instance of [`RabbitMQ`]
    ///
    /// # Examples
    /// ```
    /// use common::broker::rabbitmq::{AddressParams, RabbitMQ};
    ///
    /// let params = AddressParams {
    ///     ..AddressParams::default(),
    ///     host: "localhost",
    ///     port: 5672,
    ///     user: "guest",
    ///     password: "guest",
    /// };
    ///
    /// RabbitMQ::new(params, "my-service-queue");
    /// ```
    pub fn new(address_params: AddressParams, exchange_name: &str) -> Self {
        let string_connection = format!(
            "amqp://{}:{}@{}:{}/{}",
            address_params.user,
            address_params.password,
            address_params.host,
            address_params.port,
            address_params.vhost
        );

        RabbitMQ {
            uri: string_connection,
            connection: None,
            channel: None,
            exchange_name: exchange_name.into(),
        }
    }

    /// Connect to the RabbitMQ server.
    ///
    /// When connecting to a RabbitMQ server, you can define
    /// your own connection name by specifying the `connection_name` argument.
    /// This will be the displayed name on the dashboard.
    ///
    /// # Arguments
    /// * `connection_name` - The name of the connection. Leave it empty for an automatically assigned connection name.
    ///
    /// # Errors
    /// The connection has been established if no error is thrown.
    /// See [struct.Connection.html](https://docs.rs/lapin/4.12.3/lapin/struct.Connection.html#method.connect)
    pub async fn connect_as(&mut self, connection_name: &str) -> lapin::Result<&mut Self> {
        let uri = self.uri.as_str();
        let properties = ConnectionProperties::default()
            .with_connection_name(connection_name.into());

        let connection = Connection::connect(uri, properties).await?;
        self.connection = Some(connection);
        self.create_channel().await?;
        self.ensure_exchange_exists(self.exchange_name.as_str()).await?;

        Ok(self)
    }

    /// Create the TCP channel that will be used to communicate with the RabbitMQ server.
    ///
    /// This method requires connection to be initiated.
    /// More information at [docs.rs/lapin/4.12.3/lapin/struct.Channel.html](https://docs.rs/lapin/4.12.3/lapin/struct.Channel.html).
    ///
    /// # Errors
    /// Throw an error [`lapin::Error`] if no connection has been made.
    async fn create_channel(&mut self) -> lapin::Result<&mut Self> {
        // Checks if the connection exists
        let connection = self.connection
            .as_ref()
            .ok_or(lapin::ErrorKind::InvalidConnectionState(ConnectionState::Closed))?;

        let channel = connection.create_channel().await?;
        self.channel = Some(channel);

        Ok(self)
    }

    /// Ensure that the specified queue exists with the following options :
    /// * `durable`: `true`
    /// * `exclusive`: `false`
    /// * `auto_delete`: `false`
    ///
    /// # Notes
    ///
    /// If the queue does not exist, then it is automatically created.
    ///
    /// If the queue already exists :
    /// * with different options, an error is thrown.
    /// * with the same options, nothing happen.
    ///
    /// # Arguments
    ///
    /// * `queue_name` - The name of the queue
    ///
    /// # Errors
    ///
    /// Throw an error [`lapin::Error`] if :
    /// * no connection has been made.
    /// * the queue already exists with different options.
    pub async fn ensure_queue_exists(&self, queue_name: &str) -> lapin::Result<&Self> {
        // Checks if the channel exists
        let channel = self.channel
            .as_ref()
            .ok_or(lapin::ErrorKind::InvalidChannelState(ChannelState::Closed, "ensure_queue_exists"))?;

        let queue_options = QueueDeclareOptions {
            durable: true,
            exclusive: false,
            auto_delete: false,
            ..QueueDeclareOptions::default()
        };

        channel.queue_declare(
            queue_name.into(),
            queue_options,
            FieldTable::default()
        ).await?;

        Ok(self)
    }

    /// Ensure that the specified exchange exists with the following options :
    /// * `durable`: `true`
    ///
    /// # Notes
    ///
    /// If the exchange does not exist, then it is automatically created.
    ///
    /// If the exchange already exists :
    /// * with different options, an error is thrown.
    /// * with the same options, nothing happen.
    ///
    /// # Arguments
    ///
    /// * `exchange_name` - The name of the exchange
    ///
    /// # Errors
    ///
    /// Throw an error [`lapin::Error`] if :
    /// * no connection has been made.
    /// * the exchange already exists with different options.
    async fn ensure_exchange_exists(&self, exchange_name: &str) -> lapin::Result<&Self> {
        // Checks if the channel exists
        let channel = self.channel
            .as_ref()
            .ok_or(lapin::ErrorKind::InvalidChannelState(ChannelState::Closed, "ensure_exchange_exists"))?;

        let kind = ExchangeKind::Topic;
        let options = ExchangeDeclareOptions {
            durable: true,
            ..ExchangeDeclareOptions::default()
        };

        channel.exchange_declare(
            exchange_name.into(),
            kind,
            options,
            FieldTable::default()
        ).await?;

        Ok(self)
    }

    /// Bind a queue with an event.
    ///
    /// # Notes
    ///
    /// When this event will be published, the specified queue will receive the event.
    ///
    /// # Arguments
    ///
    /// * `queue_name` - The name of the queue to bind.
    /// * `event_name` - The name of the event to bind.
    ///
    /// # Errors
    ///
    /// Throw an error [`lapin::Error`] if :
    /// * no connection has been made.
    /// * the exchange already exists with different options.
    pub async fn bind_event(&self, queue_name: &str, event_name: &str) -> lapin::Result<&Self> {
        // Checks if the channel exists
        let channel = self.channel
            .as_ref()
            .ok_or(lapin::ErrorKind::InvalidChannelState(ChannelState::Closed, "ensure_exchange_exists"))?;

        channel.queue_bind(
            queue_name.into(),
            self.exchange_name.clone(),
            event_name.into(),
            QueueBindOptions::default(),
            FieldTable::default()
        ).await?;

        Ok(self)
    }

    /// Publish a new event
    ///
    /// # Arguments
    ///
    /// * `event_name` - The name of the event to publish
    /// * `payload` - The payload to send with the event
    ///
    /// # Errors
    ///
    /// Throw an error [`lapin::Error`] if :
    /// * no connection has been made.
    /// * the exchange already exists with different options.
    pub async fn publish_event(&self, event_name: &str, payload: &[u8]) -> lapin::Result<&Self> {
        // Checks if the channel exists
        let channel = self.channel
            .as_ref()
            .ok_or(lapin::ErrorKind::InvalidChannelState(ChannelState::Closed, "publish"))?;

        channel
            .basic_publish(
                self.exchange_name.as_str().into(),
                event_name.into(),
                BasicPublishOptions::default(),
                payload,
                BasicProperties::default(),
            )
            .await?  // Publish message
            .await?; // Wait for rabbit-mq to receive the message

        Ok(self)
    }

}
