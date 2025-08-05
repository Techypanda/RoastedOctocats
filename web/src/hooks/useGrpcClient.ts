import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';
import { OctoRoasterAPIClient } from '../proto/roasted_octocat.client'

export const useGrpcClient = () => {
    const transport = new GrpcWebFetchTransport({
        baseUrl: 'http://localhost:8081'
    });
    const client = new OctoRoasterAPIClient(transport);
    return client;
}