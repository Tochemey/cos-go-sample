VERSION 0.8

FROM golang:1.25.5-alpine

# install gcc dependencies into alpine for CGO
RUN apk --no-cache add git ca-certificates gcc musl-dev libc-dev binutils-gold curl openssh

# install docker tools
# https://docs.docker.com/engine/install/debian/
RUN apk add --update --no-cache docker

# install the go generator plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN export PATH="$PATH:$(go env GOPATH)/bin"

# install buf from source
RUN GO111MODULE=on GOBIN=/usr/local/bin go install github.com/bufbuild/buf/cmd/buf@v1.63.0

# install the various tools to generate connect-go
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

# install linter
# binary will be $(go env GOPATH)/bin/golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.8.0
RUN golangci-lint --version

# install vektra/mockery
RUN go install github.com/vektra/mockery/v2@v2.53.2

test:
  #BUILD +lint
  BUILD +local-test

protogen:
	# copy the proto files to generate
	COPY --dir protos/ .
	COPY buf.yaml buf.gen.yaml ./

    # generate the pbs
    RUN buf generate \
            --template buf.gen.yaml \
            --path protos/local/accounts \
            --path protos/chief-of-state-protos/chief_of_state

    # save artifact to gen
    SAVE ARTIFACT gen gen AS LOCAL gen

code:
    WORKDIR /app

    # download deps
    COPY go.mod go.sum ./
    RUN go mod download -x

    # copy in code
    COPY --dir +protogen/gen ./
    COPY --dir app ./

vendor:
    FROM +code

    COPY +mock/mocks ./mocks

    RUN go mod tidy && go mod vendor
    SAVE ARTIFACT /app /files

lint:
    FROM +vendor

    COPY .golangci.yml ./

    RUN golangci-lint run

mock:
    # copy in the necessary files that need mock generated code
    FROM +code
    RUN mkdir ./mocks
    # generate the mocks
    RUN mockery --all --recursive --keeptree --output ./mocks --case snake

    SAVE ARTIFACT ./mocks mocks AS LOCAL mocks

local-test:
    FROM +vendor

    WITH DOCKER --pull postgres:11
        RUN go test -mod=vendor ./app/... -race -v -coverprofile=coverage.out -covermode=atomic -coverpkg=./app/...
    END

    SAVE ARTIFACT coverage.out AS LOCAL coverage.out

compile:
    WORKDIR /build

    COPY +vendor/files ./

    RUN go build -mod=vendor  -o bin/accounts ./app/main.go

    SAVE ARTIFACT bin/accounts /accounts

docker-image:
     FROM alpine:3.16.2

    # define the image version
    ARG VERSION=dev

    WORKDIR /app
    COPY +compile/accounts ./accounts
    RUN chmod +x ./accounts

    # we listen to rpc calls on this port
    EXPOSE 50051
    # we start our prometheus server on this port
    EXPOSE 9092

    ENTRYPOINT ["./accounts"]

    SAVE IMAGE --push accounts:${VERSION}
