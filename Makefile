BINARY='shimmer.wasm'
GOBIN='${HOME}/play/gosource/go/bin/go'

all: build

build:
	GOOS=js GOARCH=wasm ${GOBIN} build -o ${BINARY} -ldflags "-s -w" ./cmd/shimmer