BINARY_NAME=mfa

build:
	go build -o bin/$(BINARY_NAME) ./cmd/mfa/

build-small:
	go build -ldflags="-s -w" -trimpath -o bin/$(BINARY_NAME) ./cmd/mfa/

completions: build
	mkdir -p completions
	./bin/$(BINARY_NAME) completion zsh > completions/mfa.zsh
	./bin/$(BINARY_NAME) completion fish > completions/mfa.fish
	./bin/$(BINARY_NAME) completion bash > completions/mfa.bash

clean:
	rm -rf bin/ completions/ dist/
