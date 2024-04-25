GO     := go
DOCKER := docker

# Directories
##################################################################
DIR_DIST     := ./dist
DIR_CMD      := ./cmd
DIR_INTERNAL := ./internal

# Build targets
##################################################################
.PHONY: build
build:
	@$(GO) build -o $(DIR_DIST)/ $(DIR_CMD)/...

.PHONY: fmt
fmt:
	gofmt -s -w -e $(DIR_CMD) $(DIR_INTERNAL)
	gci write \
		-s Standard \
		-s Default \
		-s 'Prefix(github.com)' \
		-s 'Prefix(github.com/alexZaicev/message-broker)' \
		$(DIR_CMD) $(DIR_INTERNAL)
	goimports -local github.com/alexZaicev -w $(DIR_CMD) $(DIR_INTERNAL)

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
		/protobuf/proto/messagebroker/v1alpha1/*.proto
