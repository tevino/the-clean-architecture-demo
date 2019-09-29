todo: ./cmd/cui/main.go
	go build -o $@ $^ 

.PHONY: test
test:
	go generate ./...
	go test ./...
