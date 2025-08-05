package igithub

type GithubClient interface {
	Login(request GithubLoginRequest) (*GithubLoginResponse, error)
	Refresh(request GithubRefreshRequest) (*GithubLoginResponse, error)
	WhoAmI(oauthToken string) (*WhoAmIResponse, error)
	ListEvents(request GithubListEventsRequest) (*[]GithubListEventsEvent, error)
}
