run:
	go run main.go
test:
	go test -v ./...

.PHONY:mockgen
mockgen:
    mockgen -source=internal/interface.go -destination tests/mocks/interface.go
    mockgen -source=internal/utils/utils.go -destination tests/mocks/utils/utils.go
    mockgen -source=go-sdk/pkg/chain/types.go -destination tests/mocks/chain/chain.go