ifndef MAKE_DEBUG
    MAKEFLAGS += -s
endif

.PHONY: install.buf
install.buf:
	echo "=> Installing buff to local bin dir"
	BIN="./bin" && \
	VERSION="0.38.0" && \
	BINARY_NAME="buf" && \
	curl -sSL "https://github.com/bufbuild/buf/releases/download/v$${VERSION}/$${BINARY_NAME}-`uname -s`-`uname -m`" \
		-o "$${BIN}/$${BINARY_NAME}" && \
		chmod +x "$${BIN}/$${BINARY_NAME}"

.PHONY: build
build:
	go build -o ./bin/rsk-search .

.PHONY: install.tools
install.tools:
	for TOOL in `cat tools.go | grep _ | awk -F'"' '{print $$2}'`; do \
		go install $${TOOL}; \
	done;

.PHONY: vendor
vendor:
	echo "=> Installing go deps..."
	go mod download
	echo "=> Installing proto deps..."
	./bin/buf beta mod update

.PHONY: generate.proto
generate.proto:
	rm -rf gen_temp && ./bin/buf generate && rm -rf gen && mv gen_temp gen
