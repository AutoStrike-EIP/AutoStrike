use common::event::{AddressParams, EventManager, EventPayload};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let params = AddressParams {
        host: "127.0.0.1".into(),
        password: "guest".into(),
        user: "guest".into(),
        ..AddressParams::default()
    };
    let exchange_name = "autostrike.events";
    let inbox_name = "manager_inbox";
    let mut event_manager = EventManager::new(params, exchange_name, inbox_name);
    println!("Event manager created");

    let custom_var = String::from("test");
    event_manager.on("my-event".into(), move |payload: &EventPayload| {
        println!(
            "Message received: {}",
            String::from_utf8(payload.clone()).unwrap().to_string()
        );
        println!("data: {}", custom_var);
        Ok("".into())
    });
    println!("Event registered");

    let handler = event_manager.run_in_background("task-manager", "task-manager-consumer");
    println!("Listener started");

    handler.await.unwrap()?;

    Ok(())
}
