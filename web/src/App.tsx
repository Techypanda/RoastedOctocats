import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { GatedHealthCheck } from './components/GatedHealthCheck';
import { PageLayout } from './components/PageLayout';
import { GithubAuth } from "./components/GithubAuth";
import { GithubRoaster } from "./pages/GithubRoaster";
import { BrowserRouter, Route, Routes } from "react-router";
import { GithubRoasterResult } from "./pages/GithubRoasterResult";

const queryClient = new QueryClient();

function App() {
    return (
        <QueryClientProvider client={queryClient}>
            <GatedHealthCheck>
                <PageLayout>
                    <GithubAuth>
                        <BrowserRouter>
                            <Routes>
                                <Route path="*" element={<GithubRoaster />} />
                                <Route path="/:idempotencyToken" element={<GithubRoasterResult />} />
                            </Routes>
                        </BrowserRouter>
                    </GithubAuth>
                </PageLayout>
            </GatedHealthCheck>
        </QueryClientProvider>
    )
}

export default App
