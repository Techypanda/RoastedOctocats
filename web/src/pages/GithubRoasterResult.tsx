import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";
import { useGrpcClient } from "../hooks/useGrpcClient";
import { getToken } from "../hooks/useGithubAuth";
import type { GetParsedGithubResultResponse } from "../proto/roasted_octocat";
import { Spinner } from "@fluentui/react-components";

export function GithubRoasterResult() {
    const { idempotencyToken } = useParams();
    const client = useGrpcClient();
    const { isLoading, isError, data } = useQuery({
        queryKey: [idempotencyToken],
        retry: 3,
        queryFn: async () => {
            let response: GetParsedGithubResultResponse
            do {
                response = (await client.getParsedGithubResult({ idempotencyToken: idempotencyToken || '', githubToken: getToken()?.accessToken || '' })).response;
                if (response.status != "inProgress") {
                    return response;
                }
                await new Promise(r => setTimeout(r, 500));
            } while (response.status == "inProgress");
        }
    })
    console.log(isError);
    return (
        <>
            {isLoading ? <Spinner /> : data?.result}
        </>
    )
}