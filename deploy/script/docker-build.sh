#!/usr/bin/env sh

GIT_TAG=`git describe --tags --exact-match 2>/dev/null || echo "unknown"`
DOCKER_NAMESPACE="warmans"

if [ -z "${DOCKER_IMAGE_NAME}" ]; then echo "DOCKER_IMAGE_NAME not set" && exit 1; fi

docker build -t ${DOCKER_NAMESPACE}/${DOCKER_IMAGE_NAME}:${GIT_TAG} -t ${DOCKER_NAMESPACE}/${DOCKER_IMAGE_NAME}:latest . && \
 echo "Built ${DOCKER_NAMESPACE}/${DOCKER_IMAGE_NAME}:${GIT_TAG} and ${DOCKER_NAMESPACE}/${DOCKER_IMAGE_NAME}:latest"
