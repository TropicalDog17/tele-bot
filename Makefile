run:
	go run main.go
test:
	go test -v ./...

gen-mock:
	mockgen -source=internal/interface.go -destination tests/mocks/interface.go