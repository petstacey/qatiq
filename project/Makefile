BROKER_BINARY = brokerService
AUTH_BINARY = authenticationService
LOGGER_BINARY = loggerService

up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

up_build: build_broker build_auth build_logger
	@echo "Stopping docker images (if running)..."
	docker-compose down
	@echo "Building (as required) and starting docked images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

build_auth:
	@echo "Building auth binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"