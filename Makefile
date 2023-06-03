default:
	@go build -o bin/api

run: default
	@./bin/api

test:
	@go test -v ./...

seed:
	@go run scripts/seed.go

docker:
	@docker build -t api .
	@docker run -p 3000:3000 api
