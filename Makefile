GO := go

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