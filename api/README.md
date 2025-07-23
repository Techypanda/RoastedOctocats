### Problems

- CMakePresets.json has a hardcoded path to my local machine vcpkg, this is likely not what you want.

### Build Protos

```
protoc -I="./proto" --cpp_out="./proto" ./proto/*.proto
(replace with your path to grpc installation)
protoc -I="./proto" --grpc_out="./proto" --plugin=protoc-gen-grpc="C:\Users\jonwright\vcpkg\packages\grpc_x64-windows\tools\grpc\grpc_cpp_plugin.exe" ./proto/*.proto
```