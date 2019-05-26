build: 
	go build

wasm:
	GOARCH=wasm GOOS=js go build -o facet-engine.wasm .
	mv facet-engine.wasm wasm

test:
	go test -timeout 20s -race -coverprofile coverage.txt -covermode=atomic ./...
.PHONY: test

run:
	make -B wasm
	http-server wasm