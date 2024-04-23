GO     := go
DOCKER := docker

# Directories
##################################################################
DIR_DIST := ./dist
DIR_CMD  := ./cmd

# Build targets
##################################################################

.PHONY: build
build:
	@$(GO) build -o $(DIR_DIST)/ $(DIR_CMD)/...

# Run Targets
##################################################################
.PHONY: run-%
run-%:
	$(eval NAME := $(subst run-,,$@))
	@$(DIR_DIST)/$(NAME)

# Protobuffers
##################################################################

.PHONY: gen-proto
gen-proto:
	$(eval DIR := $(shell pwd))
	$(eval UID := $(shell id -u))
	$(DOCKER) run --rm \
		-u $(UID) \
		-v $(DIR)/protobuf:/protobuf \
		-w /protobuf \
		rvolosatovs/protoc:v4.1.0 \
		--proto_path=/protobuf/proto \
		--go_out=/protobuf/go --go_opt=paths=source_relative \
		--go-grpc_out=/protobuf/go --go-grpc_opt=paths=source_relative \
		/protobuf/proto/testservice/v1/service.proto
