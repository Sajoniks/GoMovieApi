FROM golang:alpine AS build-env

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY internal internal
COPY pkg pkg
COPY cmd cmd

RUN GOOS=linux go build -o app ./cmd

FROM alpine AS runtime

WORKDIR /app
COPY --from=build-env /src/app .
ENTRYPOINT ./app