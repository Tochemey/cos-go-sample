## Sample Application

[![build](https://github.com/Tochemey/cos-go-sample/actions/workflows/build.yml/badge.svg)](https://github.com/Tochemey/cos-go-sample/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/Tochemey/cos-go-sample/branch/main/graph/badge.svg?token=6gqm0NkTJf)](https://codecov.io/gh/Tochemey/cos-go-sample)

### Overview
This is to demonstrate how to build a distributed and fault-tolerant event-sourcing/cqrs application in Go using [Chief of State](https://github.com/chief-of-state/chief-of-state).
The project adheres to [Semantic Versioning](https://semver.org) and [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

### Features

#### API Requests
- [Open Account](protos/local/accounts/v1/service.proto)
- [Credit Account](protos/local/accounts/v1/service.proto)
- [Debit Account](protos/local/accounts/v1/service.proto)
- [Get Account](protos/local/accounts/v1/service.proto)

#### Commands
- [OpenAccount](protos/local/accounts/v1/commands.proto)
- [CreditAccount](protos/local/accounts/v1/commands.proto)
- [DebitAccount](protos/local/accounts/v1/commands.proto)

#### Events
- [AccountOpened](protos/local/accounts/v1/events.proto)
- [AccountCredited](protos/local/accounts/v1/events.proto)
- [AccountDebited](protos/local/accounts/v1/events.proto)

#### State
- [BankAccount](protos/local/accounts/v1/state.proto)

#### Observability
- [Tracing](docker/otel-collector.yaml)
- [Metrics](docker/prometheus.yml)

### Quickstart
```bash
# download earthly using the following command on macos
brew install earthly/earthly/earthly

# clone the repo
git clone git@github.com:Tochemey/cos-go-sample.git

# update the git submodule
git submodule update --init

# generate the protobuf binaries and mocks
earthly +protogen
earthly +mock

# update the go dependencies
go mod tidy && go mod vendor

# build image
earthly +docker-image

# starts the prometheus, jaeger and postgres container
docker-compose up -d

# starts the application container
docker-compose --profile application up     

# OTHER HELPFUL COMMANDS

# supervise app logs
docker-compose logs -f --tail="all" accounts

# view traces
http://localhost:16686
```
