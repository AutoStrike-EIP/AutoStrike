BUILDER=cargo

build-release:
	$(BUILDER) build --all --release

build-dev:
	$(BUILDER) build --all

run-manager:
	$(BUILDER) run --bin task-manager

run-transmitter:
	$(BUILDER) run --bin task-transmitter

lint:
	$(BUILDER) clippy

clean:
	$(BUILDER) clean

.PHONY: build-release build-dev run-manager run-transmitter