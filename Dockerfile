# Build
FROM golang:1.16-alpine AS build

RUN apk add --update --no-cache --update-cache ca-certificates bash make build-base
WORKDIR /api

COPY ./ ./
RUN go mod download

RUN make build

# Deploy
FROM alpine:latest

RUN apk add --update --no-cache --update-cache ca-certificates bash make build-base
WORKDIR /

COPY --from=build /api/app/api /main

EXPOSE 3000

ENTRYPOINT ["/main"]
