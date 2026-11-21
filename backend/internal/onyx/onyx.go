package onyx

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	issuer       string
	clientID     string
	clientSecret string
	redirectURI  string
	http         *http.Client
}

type Profile struct {
	ID       string `json:"sub"`
	Username string `json:"username"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func New(issuer, clientID, clientSecret, redirectURI string) *Client {
	return &Client{
		issuer:       strings.TrimRight(issuer, "/"),
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		http:         &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) Enabled() bool {
	return c != nil && c.issuer != "" && c.clientID != "" && c.clientSecret != ""
}

func RandomString() string {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(buffer)
}

func (c *Client) AuthorizeURL(state, verifier string) string {
	challenge := sha256.Sum256([]byte(verifier))
	query := url.Values{
		"response_type":         {"code"},
		"client_id":             {c.clientID},
		"redirect_uri":          {c.redirectURI},
		"scope":                 {"username"},
		"state":                 {state},
		"code_challenge":        {base64.RawURLEncoding.EncodeToString(challenge[:])},
		"code_challenge_method": {"S256"},
	}
	return c.issuer + "/oauth/authorize?" + query.Encode()
}

func (c *Client) Identify(ctx context.Context, code, verifier string) (*Profile, error) {
	tokens, err := c.exchange(ctx, code, verifier)
	if err != nil {
		return nil, err
	}
	defer c.revoke(tokens.RefreshToken)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.issuer+"/oauth/userinfo", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
	response, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("onyx userinfo: status %d", response.StatusCode)
	}
	var profile Profile
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<16)).Decode(&profile); err != nil {
		return nil, err
	}
	if profile.ID == "" || profile.Username == "" {
		return nil, errors.New("onyx userinfo: missing sub or username")
	}
	return &profile, nil
}

func (c *Client) exchange(ctx context.Context, code, verifier string) (*tokenResponse, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {c.redirectURI},
		"code_verifier": {verifier},
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.issuer+"/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(url.QueryEscape(c.clientID), url.QueryEscape(c.clientSecret))
	response, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("onyx token: status %d", response.StatusCode)
	}
	var tokens tokenResponse
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<16)).Decode(&tokens); err != nil {
		return nil, err
	}
	if tokens.AccessToken == "" {
		return nil, errors.New("onyx token: no access token")
	}
	return &tokens, nil
}

func (c *Client) revoke(refreshToken string) {
	if refreshToken == "" {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		form := url.Values{"token": {refreshToken}}
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.issuer+"/oauth/revoke", strings.NewReader(form.Encode()))
		if err != nil {
			return
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.SetBasicAuth(url.QueryEscape(c.clientID), url.QueryEscape(c.clientSecret))
		if response, err := c.http.Do(request); err == nil {
			response.Body.Close()
		}
	}()
}
