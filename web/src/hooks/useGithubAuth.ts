import { useState } from "react";
import { useGrpcClient } from "./useGrpcClient";
import type { OctoRoasterAPIClient } from "../proto/roasted_octocat.client";

export interface GithubTokenStorage {
    accessToken: string;
    accessTokenExpiry: number;
    refreshToken: string;
    refreshTokenExpiry: number;
}

export const clearToken = () => sessionStorage.removeItem('token');
export const updateToken = (s: GithubTokenStorage) => sessionStorage.setItem('token', JSON.stringify(s));
export const getToken = (): GithubTokenStorage | null => {
    try {
        return JSON.parse(sessionStorage.getItem('token')!) as GithubTokenStorage;
    } catch {
        return null;
    }
}

const refreshAccessTokenOnExpiry = (client: OctoRoasterAPIClient, setToken: React.Dispatch<React.SetStateAction<GithubTokenStorage | null>>) => async () => {
    const token = getToken();
    try {
        const { response } = await client.refresh({
            clientId: "Iv23ligun1uyOZYdvxnq",
            refreshToken: token!.refreshToken
        });
        const accessTokenExpiryTime = (new Date().valueOf()) + (response.accessTokenExpiry * 1000);
        const refreshTokenExpiryTime = (new Date().valueOf()) + (response.refreshTokenExpiry * 1000);
        const newToken: GithubTokenStorage = {
            accessToken: response.accessToken,
            accessTokenExpiry: accessTokenExpiryTime,
            refreshToken: response.refreshToken,
            refreshTokenExpiry: refreshTokenExpiryTime
        }
        updateToken(newToken);
        setToken(newToken);
    } catch (err) {
        clearToken();
        window.location.reload();
    }
}

const clearTokensOnRefreshExpiry = () => {
    clearToken();
    window.location.reload();
}

export const useGithubAuth = () => {
    const [token, setToken] = useState(getToken());
    const client = useGrpcClient();
    if (token) {
        setTimeout(refreshAccessTokenOnExpiry(client, setToken), token.accessTokenExpiry);
        setTimeout(clearTokensOnRefreshExpiry, token.refreshTokenExpiry);
    }
    return { hasToken: !!token };
}