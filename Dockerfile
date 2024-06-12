FROM golang:alpine AS build-env

RUN apk add curl zip
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY internal internal
COPY proto proto
COPY cmd cmd

RUN curl -LO "https://github.com/protocolbuffers/protobuf/releases/download/v27.1/protoc-27.1-linux-x86_64.zip" && \
    unzip protoc-27.1-linux-x86_64.zip -d $HOME/protoc

RUN export PATH="$PATH:$HOME/protoc/bin" && \
    export PATH="$PATH:$(go env GOPATH)/bin" && \
    protoc -I . --go_out=. --go_opt=module=sajoniks.github.io/movieApi --go-grpc_out=. --go-grpc_opt=module=sajoniks.github.io/movieApi ./proto/*.proto && \
    GOOS=linux go build -o app ./cmd

FROM alpine AS runtime

WORKDIR /app
COPY --from=build-env /src/app .
ENTRYPOINT ./app