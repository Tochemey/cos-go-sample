VERSION 0.7

FROM tochemey/docker-go:1.23.2-5.0.0

test:
  BUILD +lint
  BUILD +local-test

protogen:
	# copy the proto files to generate
	COPY --dir protos/ .
	COPY buf.work.yaml buf.gen.yaml .

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
