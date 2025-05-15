ifndef MAKE_DEBUG
    MAKEFLAGS += -s
endif

LOCAL_BIN:="$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))/.env/bin"

GIT_TAG := $(shell git describe --tags --exact-match 2>/dev/null || echo "unknown")

DEV_DB_USER ?= scrimpton
DEV_DB_PASS ?= insecure
DEV_DB_HOST ?= localhost
DEV_DB_DSN ?= "postgres://$(DEV_DB_USER):$(DEV_DB_PASS)@$(DEV_DB_HOST):5432/$(DEV_DB_USER)?sslmode=disable"

#----------------------------------------------------------------------------------------------
# Dev env setup
#----------------------------------------------------------------------------------------------

.PHONY: setup
setup: check.env install.buf install.golangci dev.generate.certs
	echo "=> Setup OK"

.PHONY: install.buf
install.buf:
	echo "=> Installing buf to $(LOCAL_BIN)"
	BIN="$(LOCAL_BIN)" && \
	VERSION="1.53.0" && \
	BINARY_NAME="buf" && \
	curl -sSL "https://github.com/bufbuild/buf/releases/download/v$${VERSION}/$${BINARY_NAME}-`uname -s`-`uname -m`" \
		-o "$${BIN}/$${BINARY_NAME}" && \
		chmod +x "$${BIN}/$${BINARY_NAME}"

.PHONY: install.golangci
install.golangci:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCAL_BIN) v2.1.2

.PHONY: check.env
check.env:
ifeq ($(shell which direnv),)
	$(error "direnv not found in PATH. This tool is required to modify the project env see: https://direnv.net/docs/installation.html (or just apt-get install direnv)")
endif
ifeq ($(shell which mp3info),)
	$(error "mp3info not found in PATH. This tool is required for some commands. try apt-get install mp3info")
endif
ifeq ($(shell which ffmpeg),)
	$(error "ffmpeg not found in PATH. This tool is required for some commands: try apt-get install ffmpeg")
endif

#----------------------------------------------------------------------------------------------
# Linting
#----------------------------------------------------------------------------------------------

.PHONT: lint
lint: lint.go lint.proto

.PHONT: lint.go
lint.go:
	golangci-lint run

.PHONT: lint.proto
lint.proto:
	buf lint

#----------------------------------------------------------------------------------------------
# build/test
#----------------------------------------------------------------------------------------------

.PHONY: vendor
vendor:
	echo "=> Installing go deps..."
	go mod download
	echo "=> Installing proto deps..."
	.env/bin/buf mod update

.PHONY: build
build:
	echo "Building..."
	go build -o ./bin/rsk-search .

.PHONY: test
test:
	go test -v ./...

#----------------------------------------------------------------------------------------------
# Generated artifacts/code
#----------------------------------------------------------------------------------------------

.PHONY: generate
generate: generate.proto generate.swagger-bundle

.PHONY: generate.proto
generate.proto:
	rm -rf gen_temp && buf generate && rm -rf gen && mv gen_temp gen

.PHONY: generate.swagger-bundle
generate.swagger-bundle:
	go run ./script/swagger-merge/main.go \
		-out ../gui/src/assets/swagger.json \
		-base ./script/swagger-merge/base.swagger.json \
		-merge "$(shell ls --format=commas ./gen/api/open-api/*)"

.PHONY: generate.artifacts
generate.artifacts: generate.plaintext generate.zip

.PHONY: generate.plaintext
generate.plaintext:
	mkdir -p ./var/gen/plaintext && ./bin/rsk-search data dump-plaintext;

.PHONY: generate.zip
generate.zip:  generate.json.zip generate.plaintext.zip

.PHONY: generate.json.zip
generate.json.zip:
	cd var/data/episodes; zip -r ../../gen/episodes-json.zip *.json

.PHONY: generate.plaintext.zip
generate.plaintext.zip:
	cd var/gen/plaintext; zip -r ../../gen/episodes-plaintext.zip *.txt

#----------------------------------------------------------------------------------------------
# Run services
#----------------------------------------------------------------------------------------------

.PHONY: init.all
init.all:
	$(MAKE) init.db init.index

.PHONY: init.db
init.db:
	rm ./var/gen/ro.sqlite3; ./bin/rsk-search db load-ro && ./bin/rsk-search db load-changelogs

.PHONY: init.index
init.index:
	rm -rf ./var/gen/rsk.bluge; ./bin/rsk-search data populate-bluge-index

.PHONY: init.transcripts
init.transcripts:
	./bin/rsk-search data init

.PHONY: run
run:
	DEBUG=true ./bin/rsk-search server \
		--rw-db-dsn=${DEV_DB_DSN} \
		--oauth-reddit-secret=${REDDIT_SECRET} \
		--oauth-discord-secret=${DISCORD_SECRET} \
		--pledge-secret=${PLEDGE_SECRET} \
		--pledge-def-email=${PLEDGE_EMAIL} \
		--pledge-def-firstname=${PLEDGE_FNAME} \
		--pledge-def-lastname=${PLEDGE_LNAME} \
		--rewards-disabled=false \
		--import-keep-files=true \
		--coffee-access-token=${COFFEE_ACCESS_TOKEN} \
		--assembly-ai-access-token=${ASSEMBLY_AI_ACCESS_TOKEN} \
		--media-base-path=${MEDIA_BASE_PATH}

.PHONY: run.discord-bot
run.discord-bot:
	DEBUG=true ./bin/rsk-search discord-bot

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
	$(MAKE) dev.refresh-data

.PHONY: dev.extract.ratings
dev.extract.ratings: build
	./bin/rsk-search db merge-transcript-ratings --rw-db-dsn=${DEV_DB_DSN} --dry-run=false

#----------------------------------------------------------------------------------------------
# Misc.
#----------------------------------------------------------------------------------------------

.PHONY: extract.prod
extract.prod: build
ifndef PROD_DB_DSN
	$(error "PROD_DB_DSN required")
endif
	./bin/rsk-search db extract-tscript --rw-db-dsn=${PROD_DB_DSN} --dry-run=false
	./bin/rsk-search db merge-transcript-changes --rw-db-dsn=${PROD_DB_DSN} --dry-run=false
	./bin/rsk-search db merge-transcript-ratings --rw-db-dsn=${PROD_DB_DSN} --dry-run=false
	./bin/rsk-search db merge-transcript-tags --rw-db-dsn=${PROD_DB_DSN} --dry-run=false
	$(MAKE) dev.infer-offsets dev.refresh-data

.PHONY: extract.dev
extract.dev: build
	./bin/rsk-search db extract-tscript --rw-db-dsn=${DEV_DB_DSN} --dry-run=false
	./bin/rsk-search db merge-transcript-changes --rw-db-dsn=${DEV_DB_DSN} --dry-run=false
	./bin/rsk-search db merge-transcript-ratings --rw-db-dsn=${DEV_DB_DSN} --dry-run=false
	./bin/rsk-search db merge-transcript-tags --rw-db-dsn=${DEV_DB_DSN} --dry-run=false
	$(MAKE) dev.infer-offsets dev.refresh-data


.PHONY: report.prod
report.prod: build
ifndef PROD_DB_DSN
	$(error "PROD_DB_DSN required")
endif
	./bin/rsk-search report reddit --rw-db-dsn=${PROD_DB_DSN}

.PHONY: dev.infer-offsets
dev.infer-offsets:
	./bin/rsk-search data infer-missing-offsets

.PHONY: dev.refresh-data
dev.refresh-data:
	./bin/rsk-search data refresh && \
 	./bin/rsk-search data infer-missing-offsets && \
	$(MAKE) update.spotify-links

.PHONY: dev.generate.certs
dev.generate.certs:
	echo "Generating TLS certificates...";
	# note: that it uses a DNS entry of "server" which is what the server is called in docker-compose.
    # if docker DNS uses a different name, it would need to be changed in here.
	mkdir -p x509 && \
	openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 \
      -nodes -keyout ./x509/server_key.pem -out ./x509/server_cert.pem -subj "/CN=scrimpton.com" \
      -addext "subjectAltName=DNS:scrimpton.com,DNS:*.scrimpton.com,DNS:server,IP:0.0.0.0,IP:127.0.0.1"

.PHONY: update.spotify-links
update.spotify-links:
ifndef SPOTIFY_CLIENT_ID
	$(error "SPOTIFY_CLIENT_ID was not defined in environment")
endif
ifndef SPOTIFY_CLIENT_SECRET
	$(error "SPOTIFY_CLIENT_SECRET was not defined in environment")
endif
	./bin/rsk-search data spotifize

#----------------------------------------------------------------------------------------------
# Packaging
#----------------------------------------------------------------------------------------------

.PHONY: docker.build
docker.build:
	DOCKER_IMAGE_NAME=rsk-search-api ../deploy/script/docker-build.sh

.PHONY: docker.run
docker.run:
	docker run -it -p 8888:8888 warmans/rsk-search-api:latest

.PHONY: docker.publish
docker.publish:
	DOCKER_IMAGE_NAME=rsk-search-api ../deploy/script/docker-publish.sh

.PHONY: release
release:
	$(MAKE) build docker.build docker.publish
