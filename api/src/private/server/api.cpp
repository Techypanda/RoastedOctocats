#include <iostream>
#include "grpc/server.h"

static constexpr std::string_view SERVER_VERSION = "0.0.1";
int main() {
	return roastedoctocat::RunServer(SERVER_VERSION);
}