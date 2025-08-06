import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';
import { OctoRoasterAPIClient } from '../proto/roasted_octocat.client'

export const useGrpcClient = () => {
    const transport = new GrpcWebFetchTransport({
        baseUrl: import.meta.env.MODE == 'development' ? 'http://localhost:8081' : 'https://octocatroastercontainerapp.greencoast-ff6e3bb5.australiasoutheast.azurecontainerapps.io'
    });
    const client = new OctoRoasterAPIClient(transport);
    return client;
}