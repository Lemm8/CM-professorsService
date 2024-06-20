build:
	@go build -o bin/cm-professorsservice

run: build
	@./bin/cm-professorsservice

test:
	@go test ./...