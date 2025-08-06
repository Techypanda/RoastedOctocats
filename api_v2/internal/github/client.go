package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"techytechster.com/roastedoctocats/internal/config"
	"techytechster.com/roastedoctocats/pkg/igithub"
)

type githubClient struct {
	client *http.Client
}

func (g *githubClient) clientSecret() string {
	return config.New().ReadConfiguration(os.Getenv("APPLICATION_ENVIRONMENT")).GetGithubSecret()
}

func (g *githubClient) githubEndpoint(path string) string {
	return fmt.Sprintf("https://github.com/%s", path)
}

func (g *githubClient) apiGithubEndpoint(path string) string {
	return fmt.Sprintf("https://api.github.com/%s", path)
}

func (g *githubClient) defaultHeadersURLEncoded() http.Header {
	return http.Header{
		"Accept":               []string{"application/json"},
		"Content-Type":         []string{"application/x-www-form-urlencoded"},
		"X-Github-Api-Version": []string{"2022-11-28"},
		"User-Agent":           []string{"Roasted-Octocats"},
	}
}

func (g *githubClient) defaultHeadersJSON() http.Header {
	return http.Header{
		"Accept":               []string{"application/vnd.github+json"},
		"Content-Type":         []string{"application/json"},
		"X-Github-Api-Version": []string{"2022-11-28"},
		"User-Agent":           []string{"Roasted-Octocats"},
	}
}

var errLoginFailed = errors.New("failed to login to github")

func (g *githubClient) oauthLogin(data url.Values) (*igithub.GithubLoginResponse, error) {
	req, err := http.NewRequest("POST", g.githubEndpoint("login/oauth/access_token"), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Join(errLoginFailed, err)
	}
	req.Header = g.defaultHeadersURLEncoded()
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, errors.Join(errLoginFailed, err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(errLoginFailed, err)
	}
	var responseParsed map[string]interface{}
	err = json.Unmarshal(respBytes, &responseParsed)
	if err != nil {
		return nil, errors.Join(errLoginFailed, err)
	}
	if errorText, present := responseParsed["error"]; present {
		slog.Error("failed oauth", "errorText", errorText, "response", responseParsed)
		return nil, errors.Join(errLoginFailed, err)
	}
	return &igithub.GithubLoginResponse{
		AccessToken:        responseParsed["access_token"].(string),
		RefreshToken:       responseParsed["refresh_token"].(string),
		AccessTokenExpiry:  int(responseParsed["expires_in"].(float64)),
		RefreshTokenExpiry: int(responseParsed["refresh_token_expires_in"].(float64)),
	}, nil
}

var errNot200ListEvents = errors.New("not 200 response for listevents")
var errFailedToReadBytesListEvents = errors.New("failed to read bytes for listevents")

func (g *githubClient) ListEvents(request igithub.GithubListEventsRequest) (*[]igithub.GithubListEventsEvent, error) {
	pageNumber := 1
	if request.Page != nil {
		pageNumber = *request.Page
	}
	pageSize := 30
	if request.PageSize != nil {
		pageSize = *request.PageSize
	}
	req, err := http.NewRequest("GET", g.apiGithubEndpoint(fmt.Sprintf("users/%s/events?per_page=%d&page=%d", request.Username, pageSize, pageNumber)), nil)
	if err != nil {
		return nil, errors.Join(err, errNot200ListEvents)
	}
	req.Header = g.defaultHeadersJSON()
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.OAuthToken))
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, errors.Join(err, errNot200ListEvents)
	}
	if resp.StatusCode != 200 {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Join(err, errNot200ListEvents)
		}
		return nil, errors.Join(err, errors.New(string(respBytes)), errNot200ListEvents)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(err, errFailedToReadBytesListEvents)
	}
	var responseParsed []igithub.GithubListEventsEvent
	err = json.Unmarshal(respBytes, &responseParsed)
	if err != nil {
		return nil, errors.Join(errFailedToReadBytesListEvents, err)
	}
	// !not in production!
	slog.Debug("listevents response", "response", responseParsed)
	return &responseParsed, nil
}

func (g *githubClient) WhoAmI(oauthToken string) (*igithub.WhoAmIResponse, error) {
	req, err := http.NewRequest("GET", g.apiGithubEndpoint("user"), nil)
	if err != nil {
		return nil, errors.Join(errLoginFailed, err)
	}
	req.Header = g.defaultHeadersJSON()
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", oauthToken))
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, errors.Join(errLoginFailed, err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(errLoginFailed, err)
	}
	if resp.StatusCode != 200 {
		slog.Error("failed login", "resp", resp)
		return nil, errors.Join(errLoginFailed)
	}
	var responseParsed map[string]interface{}
	err = json.Unmarshal(respBytes, &responseParsed)
	if err != nil {
		return nil, errors.Join(errLoginFailed, err)
	}
	// !not in production!
	slog.Debug("whoami response", "response", responseParsed)
	if errorText, present := responseParsed["error"]; present {
		slog.Error("failed whoami", "errorText", errorText, "response", responseParsed)
		return nil, errors.Join(errLoginFailed, err)
	}
	// What we need: login (username) and bio
	return &igithub.WhoAmIResponse{
		Username: responseParsed["login"].(string),
		Bio:      responseParsed["bio"].(string),
	}, nil
}

func (g *githubClient) Refresh(request igithub.GithubRefreshRequest) (*igithub.GithubLoginResponse, error) {
	data := url.Values{
		"client_id":     []string{request.ClientId},
		"client_secret": []string{g.clientSecret()},
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{request.RefreshToken},
	}
	return g.oauthLogin(data)
}

func (g *githubClient) Login(request igithub.GithubLoginRequest) (*igithub.GithubLoginResponse, error) {
	data := url.Values{
		"client_id":     []string{request.ClientId},
		"client_secret": []string{g.clientSecret()},
		"code":          []string{request.Code},
		"redirect_uri":  []string{request.RedirectURI},
		"code_verifier": []string{request.CodeChallenge},
	}
	return g.oauthLogin(data)
}

func New() igithub.GithubClient {
	return &githubClient{client: &http.Client{Timeout: 30 * time.Second}}
}
