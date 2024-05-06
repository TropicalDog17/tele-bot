run:
	go run main.go
test:
	go test -v ./...

mockgen:
	mockgen -source=internal/interface.go -destination tests/mocks/interface.go
	mockgen -source=internal/utils/utils.go -destination tests/mocks/utils/utils.go
