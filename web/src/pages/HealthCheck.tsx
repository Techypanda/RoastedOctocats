import { Spinner, makeStyles } from "@fluentui/react-components";
import { useQuery } from "@tanstack/react-query";
const useStyles = makeStyles({
    div: {
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
        width: '100vw'
    }
});
export const HealthCheck = () => {
    const styles = useStyles();
    const serverIsUpQuery = useQuery({
        queryFn: async () => {
            // TODO
        }
    })
    return (
        <div className={styles.div}>
            <Spinner labelPosition="after" label="Hold on! Checking the server is up..." />
        </div>
    );
}