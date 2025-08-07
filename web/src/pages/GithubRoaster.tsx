import { useQuery } from "@tanstack/react-query";
import { useGrpcClient } from "../hooks/useGrpcClient";
import { getToken } from "../hooks/useGithubAuth";
import { ModelPromptType } from '../proto/roasted_octocat';
import { v7 } from "uuid";
import { useNavigate } from "react-router";
import { makeStyles, Button, Select, Divider } from "@fluentui/react-components";
import { useState } from "react";

const useStyle = makeStyles({
    container: {
        padding: '0 1rem 1rem 1rem'
    },
    button: {
        marginTop: '0.5rem'
    },
    flex: {
        display: 'flex',
        alignItems: 'end'
    },
    spacer: {
        margin: '0.5rem 0'
    },
    roastStyleLabel: {
        fontWeight: '600'
    },
    select: {
        minWidth: '125px',
        marginRight: '1rem'
    }
})

export const GithubRoaster = () => {
    const [modelPrompt, setModelPrompt] = useState<ModelPromptType>(ModelPromptType.ModelPromptType_EARLY2000s);
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
            promptType: modelPrompt,
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
            <Divider className={style.spacer} />
            <div className={style.flex}>
                <div className={style.select}>
                    <label className={style.roastStyleLabel} htmlFor="select-prompttype">Roast Style</label>
                    <Select id="select-prompttype" onChange={(_, data) => setModelPrompt(Number.parseInt(data.value))}>
                        <option value={ModelPromptType.ModelPromptType_EARLY2000s}>Early 2000s</option>
                        <option value={ModelPromptType.ModelPromptType_NERD}>Nerdy</option>
                        <option value={ModelPromptType.ModelPromptType_NICE}>Be Nice!</option>
                        <option value={ModelPromptType.ModelPromptType_OLDENGLISH}>Ye Olde English</option>
                        <option value={ModelPromptType.ModelPromptType_REGINAGEORGE}>Regina George</option>
                        <option value={ModelPromptType.ModelPromptType_DISCORDMOD}>Discord Mod</option>
                        <option value={ModelPromptType.ModelPromptType_DCVILLIAN}>DC Villian</option>
                        <option value={ModelPromptType.ModelPromptType_UWUIFIED}>Uwuified...</option>
                    </Select>
                </div>
                <div>
                    <Button className={style.button} onClick={() => doRoast()}>Roast Github</Button>
                </div>
            </div>
        </div>
    )
}