#!/usr/bin/env sh

GIT_TAG=`git describe --tags --exact-match 2>/dev/null || echo "unknown"`
DOCKER_NAMESPACE="warmans"

if [ -z "${DOCKER_IMAGE_NAME}" ]; then echo "DOCKER_IMAGE_NAME not set" && exit 1; fi

docker push ${DOCKER_NAMESPACE}/${DOCKER_IMAGE_NAME}:${GIT_TAG}
