ifndef MAKE_DEBUG
    MAKEFLAGS += -s
endif

GIT_TAG := $(shell git describe --tags --exact-match 2>/dev/null || echo "unknown")

DEV_DB_USER ?= scrimpton
DEV_DB_PASS ?= insecure
DEV_DB_HOST ?= localhost
DEV_DB_DSN ?= "postgres://$(DEV_DB_USER):$(DEV_DB_PASS)@$(DEV_DB_HOST):5432/$(DEV_DB_USER)?sslmode=disable"

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

.PHONY: generate
generate: generate.proto generate.swagger-bundle

.PHONY: generate.proto
generate.proto:
	rm -rf gen_temp && ./bin/buf generate && rm -rf gen && mv gen_temp gen

.PHONY: generate.swagger-bundle
generate.swagger-bundle:
	go run ./script/swagger-merge/main.go \
		-out ../gui/src/assets/swagger.json \
		-base ./script/swagger-merge/base.swagger.json \
		-merge ./gen/api/open-api/oauth.swagger.json,./gen/api/open-api/search.swagger.json,./gen/api/open-api/contrib.swagger.json

.PHONY: generate.zip
generate.zip:
	cd var/data/episodes; zip -r ../../episodes-json.zip *.json

.PHONY: run
run: build
	DEBUG=true ./bin/rsk-search server \
		--rw-db-dsn=${DEV_DB_DSN} \
		--oauth-secret=${REDDIT_SECRET} \
		--pledge-secret=${PLEDGE_SECRET} \
		--pledge-def-email=${PLEDGE_EMAIL} \
		--pledge-def-firstname=${PLEDGE_FNAME} \
		--pledge-def-lastname=${PLEDGE_LNAME} \
		--rewards-disabled=false

.PHONY: dev.services.start
dev.services.start:
	docker-compose -f dev/docker-compose.yaml up

.PHONY: dev.services.clean
dev.services.clean:
	docker-compose -f dev/docker-compose.yaml down -v

.PHONY: dev.populate.chunks
dev.populate.chunks: build
	./bin/rsk-search db load-tscript --rw-db-dsn=${DEV_DB_DSN}

.PHONY: dev.populate.test-data
dev.populate.test-data: build
	./bin/rsk-search db create-rw-testdata --rw-db-dsn=${DEV_DB_DSN}

.PHONY: dev.extract.test-data
dev.extract.test-data: build
	./bin/rsk-search db extract-tscript --rw-db-dsn=${DEV_DB_DSN} --dry-run=true

.PHONY: update.transcriptions
update.transcriptions:
ifndef SPOTIFY_TOKEN
	$(error "SPOTIFY TOKEN was not defined")
endif
	./bin/rsk-search data transcribe && \
 	SPOTIFY_TOKEN=${SPOTIFY_TOKEN} ./bin/rsk-search data spotifize

.PHONY: init.all
init.all:
	$(MAKE) init.db init.bleve

.PHONY: init.db
init.db:
	rm ./var/ro.sqlite3; ./bin/rsk-search db load-ro

.PHONY: init.index
init.index:
	rm -rf ./var/rsk.bluge; ./bin/rsk-search bluge load

.PHONY: init.bleve
init.bleve:
	rm -rf ./var/rsk.bleve; ./bin/rsk-search bleve load


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
