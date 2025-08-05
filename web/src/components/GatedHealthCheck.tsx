import { Spinner, makeStyles, MessageBar, MessageBarBody, MessageBarTitle } from "@fluentui/react-components";
import { useQuery } from "@tanstack/react-query";
import { useGrpcClient } from "../hooks/useGrpcClient";
import { v7 } from "uuid";
import type { ReactNode } from "react";
const useStyles = makeStyles({
    div: {
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
        width: '100vw'
    }
});

const ErrorInHealthCheck = () => {
    return (
        <MessageBar intent={'error'}>
            <MessageBarBody>
                <MessageBarTitle>Failure</MessageBarTitle>
                Appologies, it seem's the server is having a hiccup, Please try again later.
            </MessageBarBody>
        </MessageBar>
    )
}

export const GatedHealthCheck = (props: { children: ReactNode }): JSX.Element => {
    const client = useGrpcClient();
    const { data, isLoading, isError } = useQuery({
        queryKey: ['healthCheck'],
        retry: 0,
        queryFn: async () => {
            const { response } = await client.ping({ idempotencyToken: v7() });
            return { response };
        }
    });
    const styles = useStyles();
    return (
        data?.response?.message == 'pong' ? <>{props.children}</> :
        <div className={data?.response?.message !== 'pong' ? styles.div : undefined} >
            {isLoading ? <Spinner labelPosition="after" label="Hold on! Checking the server is up..." /> : <>
                {isError ? <ErrorInHealthCheck /> :  <ErrorInHealthCheck />}</>}
        </div>
    );
}