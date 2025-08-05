import { useEffect, type ReactNode } from "react"
import { BrowserRouter, Route, Routes, useLocation } from "react-router";
import { GithubLoginButton } from "./GithubLoginBtn";
import { makeStyles, MessageBar, MessageBarBody, MessageBarTitle, Spinner } from "@fluentui/react-components";
import { useGrpcClient } from "../hooks/useGrpcClient";
import type { OctoRoasterAPIClient } from "../proto/roasted_octocat.client";
import { clearToken, getToken, updateToken, useGithubAuth, type GithubTokenStorage } from "../hooks/useGithubAuth";
import { useQuery } from "@tanstack/react-query";

const useGhStyle = makeStyles({
});
const RequiresGithubNotification = () => {
    const style = useGhStyle();
    return (
        <>
            <MessageBar intent={'warning'}>
                <MessageBarBody>
                    <MessageBarTitle>You need a account!</MessageBarTitle>
                    Unfortunately, the API currently requires authentication, please login with github to proceed
                </MessageBarBody>
            </MessageBar>
            <GithubLoginButton />
        </>
    )
}

const useStyle = makeStyles({
    div: {
        padding: '1rem'
    }
})

const exchangeForToken = async (client: OctoRoasterAPIClient) => {
    const url = new URL(window.location.href);
    const code = url.searchParams.get('code');
    const state = url.searchParams.get('state');
    if (!state || !code) {
        window.location.href = '/';
        return;
    }
    try {
        const json = JSON.parse(sessionStorage.getItem(state)!);
        if (!json.original) {
            window.location.href = '/';
            return;
        }
        const { response } = await client.oAuth({
            clientId: "Iv23ligun1uyOZYdvxnq",
            code,
            redirectUri: "http://localhost:52986/postlogin",
            codeChallenge: json.original
        });
        const accessTokenExpiryTime = (new Date().valueOf()) + (response.accessTokenExpiry * 1000);
        const refreshTokenExpiryTime = (new Date().valueOf()) + (response.refreshTokenExpiry * 1000);
        const tokenParsed: GithubTokenStorage = {
            accessToken: response.accessToken,
            accessTokenExpiry: accessTokenExpiryTime,
            refreshToken: response.refreshToken,
            refreshTokenExpiry: refreshTokenExpiryTime
        }
        updateToken(tokenParsed);
        window.location.href = '/';
    } catch (e) {
        window.location.href = '/';
        return;
    }
}

const AccessTokenExchange = () => {
    const client = useGrpcClient();
    useEffect(() => {
        exchangeForToken(client);
    }, []);
    return (
        <Spinner appearance="primary" label="Exchanging For Token..." />
    )
}

const UserProvider = ({ children }: { children: ReactNode }) => {
    const client = useGrpcClient();
    const { isLoading, isError } = useQuery({
        queryKey: ['user'],
        retry: 3,
        queryFn: async () => {
            const { response } = await client.whoAmI({ githubToken: getToken()?.accessToken || '' });
            return { response };
        }
    });
    if (isError) {
        clearToken();
        window.location.href = '/';
    }
    return (
        <>{isLoading ? <Spinner appearance="primary" label="Querying Github For User..." /> : <>{children}</>}</>
    )
}

export const GithubAuth = ({ children }: { children: ReactNode }) => {
    const style = useStyle();
    const { hasToken } = useGithubAuth();
    if (hasToken) {
        return <UserProvider>{children}</UserProvider>;
    }
    return (
        <div className={style.div}>
            <BrowserRouter>
                <Routes>
                    <Route path='*' element={<RequiresGithubNotification />} />
                    <Route path='/postlogin' element={<AccessTokenExchange />} />
                </Routes>
            </BrowserRouter>
        </div>
    )
}