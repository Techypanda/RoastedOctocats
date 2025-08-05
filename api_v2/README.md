### Compiling Protos

protoc --go_out=./pkg --go_opt=paths=source_relative --go-grpc_out=./pkg --go-grpc_opt=paths=source_relative proto/*.proto

#### Configuration

Modify .env to have relevant params
Modify configs/application/local.json to have relevant params


### Testing

#### Integration Test

##### Docker
```
docker compose up
```

##### Manually
```
1. Start the server => go run cmd/grpc_api/main.go
2. Run The Suite => go test -count=1 ./test/integration/...
```