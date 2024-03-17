run: build
	./bin/gotemplhtmx

install:
	go install github.com/a-h/templ/cmd/templ@latest
	go get ./...
	go mod vendor
	go mod tidy
	go mod download
	bun install -D tailwindcss
	bun install -D daisyui@latest

build:
	bunx tailwindcss -i view/css/app.css -o public/styles.css 
	templ generate view
	go build -o bin/gotemplhtmx main.go 

up: ## Database migration up
	go run cmd/migrate/main.go up

reset:
	go run cmd/reset/main.go up

down: ## Database migration down
	go run cmd/migrate/main.go down

migration: ## Migrations against the database
	migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))