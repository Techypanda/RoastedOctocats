import { useQuery } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router";
import { useGrpcClient } from "../hooks/useGrpcClient";
import { getToken } from "../hooks/useGithubAuth";
import type { GetParsedGithubResultResponse } from "../proto/roasted_octocat";
import { Button, makeStyles, MessageBar, MessageBarBody, MessageBarTitle, Spinner } from "@fluentui/react-components";

const useStyle = makeStyles({
    container: {
        padding: '0 1rem 1rem 1rem'
    },
    button: {
        marginTop: '0.5rem'
    }
})

const useResultStyle = makeStyles({
    button: {
        marginTop: '0.5rem'
    }
});
const Result = ({ data, isError }: { data: GetParsedGithubResultResponse | undefined, isError: boolean }) => {
    const style = useResultStyle();
    const navigate = useNavigate();
    if (isError) {
        return (
            <>
                <MessageBar intent={'error'}>
                    <MessageBarBody>
                        <MessageBarTitle>Failure</MessageBarTitle>
                        Appologies, it seem's the server is having a hiccup, Please try again later: {data?.error || ''}
                    </MessageBarBody>
                </MessageBar>
                <Button className={style.button} onClick={() => navigate('/')}>Go Back</Button>
            </>
        )
    }
    return (
        <>
            <div>
                {data?.result}
            </div>
            <Button className={style.button} onClick={() => navigate('/')}>Go Back</Button>
        </>
    )
}

export function GithubRoasterResult() {
    const style = useStyle();
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
    return (
        <div className={style.container}>
            {isLoading ? <Spinner label="Still Roasting Your Octocat..." /> : <Result data={data} isError={isError} />}
        </div>
    )
}