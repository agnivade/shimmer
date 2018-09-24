BINARY='shimmer.wasm'

all: build

build:
	GOOS=js GOARCH=wasm go build -o ${BINARY} -ldflags "-s -w" ./cmd/shimmer
