FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerServiceApp
MAIL_BINARY=mailApp
LISTENER_BINARY=listenerApp
FRONT_BINARY=frontEndApp

DB_MONGO=db-data/mongo
DB_POSTGRES=db-data/postgres
DB_RABBITMQ=db-data/rabbitmq
CADDY_CONFIG=caddy_config
CADDY_DATA=caddy_data

## up: starts all containers in the background without forcing build
up: db_folder
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: db_folder build_broker build_auth build_logger build_mail build_listener build_front_bin
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"

## db_folder: if db folder doesn't exist, create an empty one
db_folder:
	if [ ! -d "${DB_MONGO}" ]; then mkdir -p "${DB_MONGO}"; fi
	if [ ! -d "${DB_POSTGRES}" ]; then mkdir -p "${DB_POSTGRES}"; fi
	if [ ! -d "${DB_RABBITMQ}" ]; then mkdir -p "${DB_RABBITMQ}"; fi
	if [ ! -d "${CADDY_CONFIG}" ]; then mkdir -p "${CADDY_CONFIG}"; fi
	if [ ! -d "${CADDY_DATA}" ]; then mkdir -p "${CADDY_DATA}"; fi
	@echo "db folders checked"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"

## build_mail: builds the mail binary as a linux executable
build_mail:
	@echo "Building mail binary..."
	cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	@echo "Done!"


## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo "Building mail binary..."
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Done!"

## build_front_bin: builds the front end binary as a linux executable
build_front_bin:
	@echo "Building front end binary..."
	cd ../front-end && env GOOS=linux CGO_ENABLED=0 go build -o ${FRONT_BINARY} ./cmd/web
	@echo "Done!"

## build_front: builds the frone end binary
build_front:
	@echo "Building front end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ../front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"