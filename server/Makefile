ifndef MAKE_DEBUG
    MAKEFLAGS += -s
endif

GIT_TAG := $(shell git describe --tags --exact-match 2>/dev/null || echo "unknown")

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
install.tools: install.buf
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

.PHONY: run
run: build
	./bin/rsk-search server

.PHONY: update.transcriptions
update.transcriptions:
ifndef SPOTIFY_TOKEN
	$(error "SPOTIFY TOKEN was not defined")
endif
	# there isn't and easy way to do this because so much is built on top of the transcription.
	# even in this case any new tags will not be added for example
	./bin/rsk-search data transcribe && \
 	SPOTIFY_TOKEN=${SPOTIFY_TOKEN} ./bin/rsk-search data spotifize && \
 	./bin/rsk-search data tag-dialog

.PHONY: init.all
init.all:
	$(MAKE) init.db init.index

.PHONY: init.db
init.db:
	rm ./var/ro.sqlite3; ./bin/rsk-search db load

.PHONY: init.index
init.index:
	rm -rf ./var/rsk.bleve; ./bin/rsk-search index load

.PHONY: docker.build
docker.build:
	DOCKER_IMAGE_NAME=rsk-search-api ../deploy/script/docker-build.sh

.PHONY: docker.run
docker.run:
	docker run -it -p 8888:8888 $(DOCKER_IMAGE_NAME):latest

.PHONY: docker.publish
docker.publish:
	DOCKER_IMAGE_NAME=rsk-search-api ../deploy/script/docker-publish.sh

.PHONY: release
release:
	$(MAKE) build docker.build docker.publish
