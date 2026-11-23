/// The event's payload type: [Vec\<u8>]
pub type EventPayload = Vec<u8>;

/// The event's callback trait
///
/// # Notes
/// Implements `Send` trait for the ability to switch between threads.
/// Is defined as `'static` to ensure that the function does not include any temporary reference (&) which could be destroyed before the function is called.
pub trait EventCallbackFn: Fn(&EventPayload) -> Result<String, String> + Send + 'static {}

/// Implementation of the EventCallbackFn for all corresponding lambdas
///
/// # Notes
/// This applies [EventCallbackFn] to all lambda `T` that follows the criteria of [EventCallbackFn].
impl<T> EventCallbackFn for T where T: Fn(&EventPayload) -> Result<String, String> + Send + 'static {}
