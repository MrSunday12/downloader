# syntax=docker/dockerfile:1.3

#build stage
FROM golang:1.17.5-alpine AS build
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go get -v ./...
RUN go mod download
COPY *.go ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -v -o /go/bin/app ./...

#final stage
FROM alpine:3.15 as bin
RUN apk --no-cache add ca-certificates
COPY --from=build /go/bin/app /app
ENTRYPOINT /app
