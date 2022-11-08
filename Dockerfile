# syntax=docker/dockerfile:1.2
FROM golang:1.18.7 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./
COPY generator/sbom.go ./generator/sbom.go

RUN CGO_ENABLED=0 go build -o /sbom-generator

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /sbom-generator ./

ENTRYPOINT ["/app/sbom-generator", "-path", "/lifecycle"]
