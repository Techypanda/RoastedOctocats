#pragma once
#include <string_view>
#include "api.grpc.pb.h"
namespace roastedoctocat {
	class OctocatAPIServiceImpl final : public api::OctoRoasterAPI::Service
	{
	public:
		grpc::Status Ping(grpc::ServerContext* context, const api::PingRequest* request, api::PingResponse* response) override;
		~OctocatAPIServiceImpl() override;
		OctocatAPIServiceImpl(std::string_view ServerVersion);
	private:
		std::string_view _ServerVersion;
	};
	int RunServer(std::string_view ServerVersion, std::string_view Port = "8000", std::string_view IpAddressSourceOverride = "0.0.0.0");
}