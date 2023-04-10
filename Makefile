dep:
	@docker-compose up

test:
	@go test ./... -v -short

test.coverage:
	-go test -coverprofile cover.out -v ./...
	@go tool cover -html=cover.out

test.integration: clean development-serve
	export DB_HOST=localhost
	export DB_PORT=5432
	export DB_USER=test
	export DB_PASSWORD=password
	export DB_NAME=your-money
	go test -tags=integration -v -count=1 ./... -failfast
	make clean

load_local_env:
	@export DB_HOST=localhost
	@export DB_PORT=5432
	@export DB_USER=test
	@export DB_PASSWORD=password
	@export DB_NAME=your-money

development-serve: load_local_env
	@docker-compose up -d --build

clean:
	@docker-compose down
	@ - docker volume rm your-money_postgres_data
