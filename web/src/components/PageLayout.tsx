import { makeStyles, MessageBar, MessageBarBody, MessageBarTitle } from "@fluentui/react-components";
import type { ReactNode } from "react";

const useStyles = makeStyles({
    div: {
        padding: '1rem'
    },
    notice: {
        // marginBottom: '1rem'
    }
});

export const PageLayout = (props: { children: ReactNode }) => {
    const styles = useStyles();
    return (
        <>
        <div className={styles.div}>
            <MessageBar intent={'info'} className={styles.notice}>
                <MessageBarBody>
                    <MessageBarTitle>Notice</MessageBarTitle>
                    This application is currently under development and may not function as expected.
                </MessageBarBody>
            </MessageBar>
        </div>
        {props.children}
        </>
    )
}