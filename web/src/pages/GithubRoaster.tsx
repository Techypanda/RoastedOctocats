import { useQuery } from "@tanstack/react-query";
import { useGrpcClient } from "../hooks/useGrpcClient";
import { getToken } from "../hooks/useGithubAuth";
import { v7 } from "uuid";
import { useNavigate } from "react-router";
import { makeStyles, Button } from "@fluentui/react-components";

const useStyle = makeStyles({
    container: {
        padding: '0 1rem 1rem 1rem'
    },
    button: {
        marginTop: '0.5rem'
    }
})

export const GithubRoaster = () => {
    const style = useStyle();
    const client = useGrpcClient();
    const { data, isLoading } = useQuery({
        queryKey: ['user'],
        retry: 3,
        queryFn: async () => {
            const { response } = await client.whoAmI({ githubToken: getToken()?.accessToken || '' });
            return { response };
        }
    });
    const navigate = useNavigate();
    async function doRoast() {
        const uuid = v7();
        await client.parseGithub({
            idempotencyToken: uuid,
            githubToken: getToken()?.accessToken || ''
        });
        navigate(`/${uuid}`);
    }
    if (isLoading) {
        return <></>; // shouldn't be rendered ever
    }
    return (
        <div className={style.container}>
            <div>Hi {data?.response.username}, Click the button below to have your github reviewed</div>
            <Button className={style.button} onClick={() => doRoast()}>Roast Github</Button>
        </div>
    )
}