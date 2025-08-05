package igithub

const PushEventType string = "PushEvent"

type RepoDetails struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type CommitDetails struct {
	Message string `json:"message"`
}

type PayloadDetails struct {
	Commits []CommitDetails `json:"commits"`
}

type GithubListEventsRequest struct {
	Page       *int
	PageSize   *int
	Username   string
	OAuthToken string
}

type WhoAmIResponse struct {
	Username string
	Bio      string
}

type GithubListEventsEvent struct {
	ID      string         `json:"id"`
	Type    string         `json:"type"`
	Repo    RepoDetails    `json:"repo"`
	Payload PayloadDetails `json:"payload"`
}

type GithubLoginResponse struct {
	AccessToken        string `json:"accessToken"`
	RefreshToken       string `json:"refreshToken"`
	AccessTokenExpiry  int    `json:"accessTokenExpiry"`
	RefreshTokenExpiry int    `json:"refreshTokenExpiry"`
}

type GithubLoginRequest struct {
	ClientId      string `json:"clientId"`
	Code          string `json:"code"`
	RedirectURI   string `json:"redirectURI"`
	CodeChallenge string `json:"codeChallenge"`
}

type GithubRefreshRequest struct {
	ClientId     string `json:"clientId"`
	RefreshToken string `json:"refreshToken"`
}
