https://github.com/fullstorydev/grpcurl

go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

grpcurl -plaintext list
grpcurl -plaintext 0.0.0.0:8000 describe api.OctoRoasterAPI.Hello

# Ping Method
grpcurl -plaintext -d '{ \"idempotencyToken\": \"someToken\" }' 0.0.0.0:8000 api.OctoRoasterAPI/Ping
```
{
  "message": "pong",
  "serverVersion": "0.0.1",
  "idempotencyToken": "someToken"
}
```