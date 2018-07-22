# set default shell
SHELL := $(shell which bash)
ENV = /usr/bin/env
DOCKER_API_V1 = "seblegall/blackbeard-example-api:v1"
DOCKER_API_V2 = "seblegall/blackbeard-example-api:v2"
DOCKER_FRONT = "seblegall/blackbeard-example-front:v1"

.SHELLFLAGS = -c

.SILENT: ;               # no need for @
.ONESHELL: ;             # recipes execute in same shell
.NOTPARALLEL: ;          # wait for this target to finish
.EXPORT_ALL_VARIABLES: ; # send all vars to shell

.PHONY: all
.DEFAULT: build

build: ## Build my-app docker images
	docker build -t ${DOCKER_API_V1} ./api/v1/
	docker build -t ${DOCKER_API_V2} ./api/v2/
	docker build -t ${DOCKER_FRONT} ./front/

push: ## Push my-app docker images to hub.docker.com
	docker push ${DOCKER_API_V1}
	docker push ${DOCKER_API_V2}
	docker push ${DOCKER_FRONT}