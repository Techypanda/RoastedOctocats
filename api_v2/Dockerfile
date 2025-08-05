# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder
WORKDIR /app
COPY . /app/
RUN apt update
RUN apt install -y ca-certificates
RUN update-ca-certificates
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/grpc_api/main.go

FROM scratch AS final
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/main ./main
COPY --from=builder /app/configs /app/configs
EXPOSE 8080
ENTRYPOINT ["/app/main"]