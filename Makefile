# Builds the project under the `fwd-dog` file name
.PHONY: build
build:
	go build -o fwd-dog main.go

# Runs the app
.PHONY: run
run:
	go run main.go

# Runs the project's automated tests
.PHONY: test
test:
	go test ./...

# Runs the tests in the local environment
.PHONY: cover
cover:
	@go test -coverprofile coverage.out -coverpkg ./... ./... && go tool cover -html=coverage.out -o coverage.html && rm coverage.out

# Starts up the development environment
.PHONY: dev
dev:
	-@docker-compose -f .docker/docker-compose.yml -f .docker/docker-compose.dev.yml up --build --abort-on-container-exit --renew-anon-volumes --remove-orphans

# Executes 'bash' on the 'app' container
.PHONY: enter
enter:
	$(eval ID := $(shell docker ps -q --filter "label=app=fwd-dog"))
	@if [ -z $(ID) ]; then echo >&2 "ERROR: The 'app' container is missing" && exit 1; fi
	# Entering the container
	-@docker exec -it $(ID) bash

# Builds a temporary image and runs the tests
.PHONY: isolated-test
isolated-test:
	# Running the tests
	@docker-compose -f .docker/docker-compose.yml -f .docker/docker-compose.test.yml up --build --abort-on-container-exit --renew-anon-volumes --remove-orphans
