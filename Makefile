NAME ?= hms-trs-kafkalib
VERSION ?= $(shell cat .version)

all: image unittest

image:
	docker build --pull ${DOCKER_ARGS} --tag '${NAME}:${VERSION}' .

unittest:
	./runUnitTest.sh

