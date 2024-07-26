build:
	@go build -o cmd/api
run: build
	@./cmd/api
test:
	@go test -v ./..