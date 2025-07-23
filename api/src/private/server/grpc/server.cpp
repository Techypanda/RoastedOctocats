#include <iostream>
#include <grpcpp/grpcpp.h>
#include <grpcpp/health_check_service_interface.h>
#include <grpcpp/ext/proto_server_reflection_plugin.h>
#include "server.h"

namespace roastedoctocat {
	int RunServer(std::string_view ServerVersion, std::string_view Port, std::string_view IpAddressSourceOverride)
	{
		std::ostringstream ServerAddress;
		ServerAddress << IpAddressSourceOverride << ":" << Port;
		roastedoctocat::OctocatAPIServiceImpl service{ServerVersion};
		grpc::EnableDefaultHealthCheckService(true);
		grpc::reflection::InitProtoReflectionServerBuilderPlugin();
		grpc::ServerBuilder builder;
		builder.AddListeningPort(ServerAddress.str(), grpc::InsecureServerCredentials());
		builder.RegisterService(&service);
		std::unique_ptr<grpc::Server> server(builder.BuildAndStart());
		std::cout << "GRPC API listening on " << ServerAddress.str() << std::endl;
		server->Wait();
		return 0;
	}
	grpc::Status OctocatAPIServiceImpl::Ping(grpc::ServerContext* context, const api::PingRequest* request, api::PingResponse* response)
	{
		std::cout << "Ping requested, idempotency token: " << request->idempotencytoken() << std::endl;
		response->set_message("pong");
		response->set_idempotencytoken(request->idempotencytoken());
		response->set_serverversion(_ServerVersion);
		return grpc::Status::OK;
	}
	OctocatAPIServiceImpl::~OctocatAPIServiceImpl()
	{
	}
	OctocatAPIServiceImpl::OctocatAPIServiceImpl(std::string_view ServerVersion) : _ServerVersion(ServerVersion)
	{
	}
}