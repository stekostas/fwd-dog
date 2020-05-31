# fwd.dog

![Tests](https://github.com/stekostas/fwd-dog/workflows/Tests/badge.svg)

> An ephemeral short link generator.

## Introduction

fwd.dog is an ephemeral short link generator. You provide a target URL, set an expiration time, and you get a short link. The system stores only the provided target URL bound to the generated link. No accounts, no tracking, entirely anonymous.

The main purpose of the service is to provide a very short link to your target URL, so it can be quickly typed by hand or even memorised for easy transfer across disconnected devices. The generated links are publicly accessible for the duration of their existence and are not meant to be used as permalinks to their target URLs.

The system can hold up to 56,800,235,584 unique links at any given time. The generated hash _- something like `xY0z4` -_ consists of 1 to 6 alphanumeric characters and is recycled after expiration.

## Requirements

- Go 1.14 or later

## Recommended Tools

- Docker & Docker Compose
- GNU Make

## Installation

You can run the server natively using `go` or using Docker and Docker Compose to spin up an isolated environment.

### Using Docker Compose (Recommended)

The project comes with Docker Compose configuration files (under the `.docker/` directory) that can spin up isolated environments for development, testing, and production.

To start up the server, run the provided `make isolated-run` target.

```
# Start up the production environment
$ make isolated-run
```

By default, the server is listening on `:3000`. Visit http://localhost:3000 on your browser to use the service. To stop the environment just press <kbd>Ctrl</kbd> + <kbd>C</kbd> once to gracefully shutdown or twice to force it.

### Using Native Tools

As a prerequisite, you need a Redis server running as a service and Go version 1.14 or later installed on your machine.

The server expects two environment variables to be set:

- `APP_HOST`: Specifies the host the server will be listening on (e.g. `:3000`)
- `REDIS_ADDRESS`: Specifies the Redis server address (e.g. `localhost:6379`)
 
Use the provided `Makefile` to start up the server.

```
# Start up the server listening on ":3000" with Redis on "localhost:6379"
$ APP_HOST=":3000" REDIS_ADDRESS="localhost:6379" make run
```

Alternatively, you can start the server using `go` directly.

```
$ APP_HOST=":3000" REDIS_ADDRESS="localhost:6379" go run .
```

### Using Docker

The projet's GitHub repository contains two Docker images per version.

```
# Run the latest version "buster" image
$ docker run -it --rm -e REDIS_ADDRESS="my-redis-host:6379" docker.pkg.github.com/stekostas/fwd-dog/fwd-dog:0.1.0-buster
```

## Development

For an easy and consistent development experience, you will need to have Docker and Docker Compose installed.

Simply run `make dev` to spin up the decoupled development environment and then enter the container to run, build, and test the server.

```
# Start up the development environment
$ make dev
...

# On a new terminal, enter the app container
$ make enter

# Inside the container, run the server...
app@host:/app$ make run

# ...or run the tests
app@host:/app$ make test
```

The environment mounts your local project directory, so you can make changes to the source code on the fly. Remember to restart the server when changing the source code.

The services started by the development environment:

- `app`: The Go server exposing port `3000`.
- `redis`: The Redis server on port `6379`.
- `redis_commander`: The [Redis Commander app](https://github.com/joeferner/redis-commander) accessible on port `8081` of your localhost.

## Testing

### Using Docker Compose

To run the tests in an isolated environment, you can use the `make isolated-test` target.

```
$ make isolated-test
```

This will provision and start up a testing environment using Docker Compose. It will run the tests in the container and exit, tearing down the environment as it gracefully shuts down.

### Using `make test`

On your development environment, you can simply enter the app container and run `make test`.

To get a coverage report, run `make cover`. On tests success, it will generate a `coverage.html` file on your project root. Open it on your browser to see the full coverage report.

## Versioning

The project follows the [Semantic Versioning v2](https://semver.org/spec/v2.0.0.html) standard. For the available project versions, check the [releases page](https://github.com/stekostas/fwd-dog/releases). 

## License

Copyright (c) 2020 [Kostas Stergiannis](https://github.com/stekostas)

The project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
