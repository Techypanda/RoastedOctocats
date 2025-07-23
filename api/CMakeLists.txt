cmake_minimum_required (VERSION 3.8)
if (POLICY CMP0141)
  cmake_policy(SET CMP0141 NEW)
  set(CMAKE_MSVC_DEBUG_INFORMATION_FORMAT "$<IF:$<AND:$<C_COMPILER_ID:MSVC>,$<CXX_COMPILER_ID:MSVC>>,$<$<CONFIG:Debug,RelWithDebInfo>:EditAndContinue>,$<$<CONFIG:Debug,RelWithDebInfo>:ProgramDatabase>>")
endif()
project ("api")
find_package( gRPC CONFIG REQUIRED )
find_package( Protobuf REQUIRED )
# Proto Gen (TODO: Make work)
# add_custom_command(
#    OUTPUT COMPILED_PROTO_1
#    COMMAND protoc -I="C:\\opensource_projects\\RoastedOctocats\\api\\proto" --cpp_out="C:\\opensource_projects\\RoastedOctocats\\api\\proto" "C:\\opensource_projects\\RoastedOctocats\\api\\proto\\*.proto"
# )
#add_custom_command(
#    OUTPUT COMPILED_PROTO_2
#    COMMAND protoc -I="C:\\opensource_projects\\RoastedOctocats\\api\\proto" --grpc_out="C:\\opensource_projects\\RoastedOctocats\\api\\proto" --plugin=protoc-gen-grpc="C:\\Users\\jonwright\\vcpkg\\packages\\grpc_x64-windows\\tools\\grpc\\grpc_cpp_plugin.exe" "C:\\opensource_projects\\RoastedOctocats\\api\\proto\\*.proto"
#)
#add_custom_target(compile_protos ALL DEPENDS COMPILED_PROTO_1 COMPILED_PROTO_2)
file(GLOB PROTO_FILES "proto/*.pb.cc") # DEPENDS compile_protos)
# Code Files
file(GLOB API_SOURCE_FILES "src/private/server/*.cpp" "src/private/server/**/*.cpp")
file (GLOB API_HEADER_FILES "src/private/server/*.h" "src/private/server/**/*.h")
# Executables
add_executable(api ${API_SOURCE_FILES} ${API_HEADER_FILES} ${PROTO_FILES})
target_include_directories(api PRIVATE proto)
target_link_libraries(api PRIVATE gRPC::grpc++ gRPC::grpc++_reflection gRPC::gpr gRPC::grpc gRPC::grpc++ protobuf::libprotoc protobuf::libprotobuf protobuf::libprotobuf-lite)