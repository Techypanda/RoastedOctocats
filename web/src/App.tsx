import { MessageBar, MessageBarBody, MessageBarTitle } from '@fluentui/react-components';
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { HealthCheck } from './pages/HealthCheck';
const queryClient = new QueryClient();
function App() {
    return (
        <QueryClientProvider client={queryClient}>
            {/*<MessageBar intent={'info'}>*/}
            {/*    <MessageBarBody>*/}
            {/*        <MessageBarTitle>Notice</MessageBarTitle>*/}
            {/*        This application is currently under development and may not function as expected.*/}
            {/*    </MessageBarBody>*/}
            {/*</MessageBar>*/}
            <HealthCheck />
        </QueryClientProvider>
    )
}

export default App
