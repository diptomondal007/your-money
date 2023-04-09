dep:
	@docker-compose up

test:
	@go test ./... -v

test-coverage:
	-go test -coverprofile cover.out -v ./...
	@go tool cover -html=cover.out

load_local_env:
	@export DB_HOST=localhost
	@export DB_PORT=5432
	@export DB_USER=test
	@export DB_PASSWORD=password
	@export DB_NAME=food

load_remote_env:
	@export DB_HOST=localhost
	@export DB_PORT=5432
	@export DB_USER=test
	@export DB_PASSWORD=password
	@export DB_NAME=food

development-serve: load_local_env
	@docker-compose up --build

clean:
	@docker-compose down