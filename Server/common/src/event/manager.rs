use std::collections::HashMap;
use std::sync::Arc;

use futures_lite::StreamExt;
use lapin::options::BasicAckOptions;
use lapin::types::ShortString;
use tokio::task::JoinHandle;

use super::types::{EventCallbackFn, EventPayload};
use super::{AddressParams};
use super::rabbitmq::RabbitMQService;

pub struct EventManager {
    unknown_callback: Box<dyn EventCallbackFn>,
    callbacks: HashMap<ShortString, Box<dyn EventCallbackFn>>,
    rabbitmq_service: RabbitMQService,
    inbox: String,
}

impl EventManager {
    /// Create a new instance of [EventManager].
    ///
    /// # Arguments
    /// * `address` - An instance of [AddressParams].
    /// * `exchange_name` - The name of the exchange that will be manipulated.
    /// * `inbox_name` - The name of the RabbitMQ queue where the messages will be received.
    pub fn new(address: AddressParams, exchange_name: &str, inbox_name: &str) -> Self {
        EventManager {
            unknown_callback: Box::new(Self::undefined_callback),
            rabbitmq_service: RabbitMQService::new(address, exchange_name),
            callbacks: HashMap::default(),
            inbox: inbox_name.into(),
        }
    }

    /// Bind an event to a callback.
    ///
    /// # Arguments
    /// * `event` - The name of the routing_key from RabbitMQ
    /// * `callback` - The function to call when the event is catch.
    ///     * The function must follow the prototype `fn (payload: &EventCallbackPayload) -> Result<String, String>`.
    ///
    /// # Types
    /// * `EventCallbackPayload` - A `Vec<u8>` which corresponds to the payload.
    ///
    /// # Example
    ///
    /// **Warning** - Do not forget the `move` keyword before your callback!
    ///
    /// Code example:
    /// ```
    /// let some_data = Arc::new(String::from("My shared data")); // This shared pointer is kept in the main scope
    /// let sharable_data = Arc::clone(&some_data); // This shared pointer will be moved in the callback
    ///
    /// event_manager.on("my-event".into(), move |payload: &EventCallbackPayload| {
    ///     println!("Message received: {}", String::from_utf8(payload.clone()).unwrap().to_string());
    ///     println!("data: {}", sharable_data);
    ///     Ok("".into())
    /// });
    /// ```
    pub fn on(&mut self, event: &str, callback: impl EventCallbackFn) {
        self.callbacks.insert(event.into(), Box::new(callback));
    }

    /// Publish a new event.
    ///
    /// # Arguments
    /// * `event` - The event name to publish.
    /// * `payload` - The payload of the event.
    pub async fn publish(&mut self, event: &str, payload: EventPayload) -> lapin::Result<()> {
        self.rabbitmq_service.publish_event(event, payload.as_ref()).await?;
        Ok(())
    }

    /// Call [Self::run] in a new thread.
    ///
    /// # Returns
    /// The thread's [JoinHandle] instance.
    pub fn run_in_background(
        mut self,
        service_name: &str,
        consumer_tag: &str,
    ) -> JoinHandle<lapin::Result<()>> {
        let consumer_tag = String::from(consumer_tag);
        let service_name = String::from(service_name);
        let handler = tokio::spawn(async move { self.run(service_name, consumer_tag).await });

        handler
    }

    /// Run the broker service.
    ///
    /// # Errors
    ///
    /// Throw an error [`lapin::Error`] in different cases. See [lapin documentation](https://docs.rs/lapin/4.12.3/lapin/index.html).
    ///
    pub async fn run(&mut self, service_name: String, consumer_tag: String) -> lapin::Result<()> {
        self.rabbitmq_service.connect_as(service_name.as_str()).await?;
        self.register_events().await?;

        let mut consumer =
            self.rabbitmq_service.new_consumer(consumer_tag.as_str(), self.inbox.as_str()).await?;

        while let Some(delivery) = consumer.next().await {
            let packet = Arc::new(delivery?);
            let event = &packet.routing_key;
            let callback = match self.callbacks.get(event) {
                Some(callback) => callback,
                None => &self.unknown_callback,
            };

            let result = callback(&packet.data);
            match result {
                Ok(msg) => {
                    packet.ack(BasicAckOptions::default()).await?;
                    println!("[SUCCESS] {}", msg);
                },
                Err(err) => println!("[ERROR] {}", err),
            }
        }
        Ok(())
    }

    /// Register all events bound with a callback in RabbitMQ
    async fn register_events(&mut self) -> lapin::Result<()> {
        for (event, _callback) in self.callbacks.iter_mut() {
            self.rabbitmq_service
                .ensure_queue_exists(event.as_str())
                .await?
                .bind_event(self.inbox.as_str(), event.as_str())
                .await?;
        }
        Ok(())
    }

    /// Default callback in the case that the broker receive an event that is not in the callbacks [HashMap].
    fn undefined_callback(_payload: &EventPayload) -> Result<String, String> {
        println!("Undefined callback.");
        Ok("".into())
    }
}
