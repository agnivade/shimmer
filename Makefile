BINARY='shimmer.wasm'
all: build

build:
	GOOS=js GOARCH=wasm go build -o ${BINARY} ./cmd/shimmer

build-prod:
	GOOS=js GOARCH=wasm go build -o ${BINARY} -ldflags "-s -w" ./cmd/shimmer
