# ENV Variables
POSTGRES_PORT=5432
REDIS_PORT=6379
APP_PORT=1000
NETWORK=app-network


create-network:
	docker network inspect $(NETWORK) >NUL 2>&1 || docker network create $(NETWORK)

remove-network:
	docker network rm $(NETWORK)

build-postgres:
	docker build -t custom-postgres ./docker/postgres

run-postgres:
	docker run --name postgres-server --network $(NETWORK) -p $(POSTGRES_PORT):5432 \
	-e POSTGRES_DB=department \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_PASSWORD=P@ssw0rd \
	-d custom-postgres

remove-postgres:
	docker stop postgres-server
	docker rm postgres-server

build-redis:
	docker build -t custom-redis ./docker/redis

run-redis:
	docker run --name redis-server --network $(NETWORK) -p $(REDIS_PORT):6379 \
	-d custom-redis

remove-redis:
	docker stop redis-server
	docker rm redis-server

build-app:
	docker build -t my-go-app -f docker/app/Dockerfile .

run-app:
	docker run --name go-app --network $(NETWORK) -p $(APP_PORT):1000 \
	--env-file .env \
	--link postgres-server:postgres-server \
	--link redis-server:redis-server \
	-v cert:/app/cert \
	-v keys:/app/keys \
	-v logs:/app/logs \
	-d my-go-app

remove-app:
	docker stop go-app
	docker rm go-app

start-all: create-network build-postgres run-postgres build-redis run-redis build-app run-app

stop-all: remove-postgres remove-redis remove-app remove-network

## RUN APPLICATION
run:
	@echo -e "Running the application..."
	@dotenv -e .env -- go run ./cmd/main.go

## RUN TESTS
test:
	@echo -e "Running tests..."
	@dotenv -e .env -- go test -v ./tests/department_test.go