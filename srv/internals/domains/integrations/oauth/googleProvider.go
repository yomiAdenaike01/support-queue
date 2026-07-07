package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func (p *Provider) setProperties(deps providerDeps) {
	p.Name = deps.name
	p.ClientId = deps.clientId
	p.ClientSecret = deps.clientSecret
	p.RedirectUrl = deps.redirectUri
	p.Scopes = deps.scopes
}

type googleProvider struct {
	Provider
}

type googleTokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_token_expires_in"`
}

func (g *googleTokenResponse) getAccessToken() string {
	return g.AccessToken
}

func Ptr[V any](val V) *V {
	return &val
}

func (g *googleTokenResponse) toSaveInput() saveTokenInput {
	return saveTokenInput{
		provider:         GOOGLE,
		accessToken:      g.AccessToken,
		refreshToken:     Ptr(g.RefreshToken),
		accessExpiresAt:  time.Now().Add(time.Duration(g.ExpiresIn * int(time.Second))),
		refreshExpiresAt: Ptr(time.Now().Add(time.Duration(g.RefreshExpiresIn * int(time.Second)))),
	}
}

func (g *googleProvider) authorisationEndpoint() string {
	queryVals := url.Values{
		"client_id":     []string{g.ClientId},
		"scope":         []string{g.Scopes},
		"redirect_uri":  []string{g.RedirectUrl},
		"response_type": []string{"code"},
		"access_type":   []string{"offline"},
		"prompt":        []string{"consent"},
	}
	return fmt.Sprintf("%s?%s", g.AuthEndpoint, queryVals.Encode())
}

func (g *googleProvider) tokenExchange(code string) (ITokenResponse, error) {
	var res googleTokenResponse
	rawBody, err := json.Marshal(map[string]string{
		"code": code,
	})
	if err != nil {
		return &res, err
	}
	tknEndpoint := g.TokenEndpoint
	queryParams := url.Values{
		"grant_type":    []string{"authorization_code"},
		"code":          []string{code},
		"client_id":     []string{g.ClientId},
		"client_secret": []string{g.ClientSecret},
		"redirect_uri":  []string{g.RedirectUrl},
	}.Encode()

	response, err := http.DefaultClient.Post(fmt.Sprintf("%s?%s", tknEndpoint, queryParams), "application/x-www-form-urlencoded", bytes.NewReader(rawBody))
	if err != nil {
		return &res, err
	}
	if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
		return &res, err
	}
	return &res, nil
}

func (g *googleProvider) getName() ProviderName {
	return g.Name
}

type refreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (r refreshResponse) getAccessToken() string {
	return r.AccessToken
}

func (r refreshResponse) getExpiresAt() time.Time {
	return time.Now().Add(time.Second * time.Duration(r.ExpiresIn))
}

func (r refreshResponse) toSaveInput() saveTokenInput {
	return saveTokenInput{
		accessToken:      r.AccessToken,
		accessExpiresAt:  r.getExpiresAt(),
		refreshToken:     nil,
		refreshExpiresAt: nil,
		provider:         GOOGLE,
	}
}

func (g *googleProvider) refresh(refreshToken string) (iRefreshResponse, error) {
	response, err := http.DefaultClient.Post(g.TokenEndpoint, "application/x-www-form-urlencoded", bytes.NewBufferString(url.Values{
		"grant_type":    []string{"refresh_token"},
		"client_id":     []string{g.ClientId},
		"client_secret": []string{g.ClientSecret},
		"refresh_token": []string{refreshToken},
	}.Encode()))
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(response.Body)
	if response.StatusCode != http.StatusOK {
		var badResponse map[string]any
		if err := decoder.Decode(&badResponse); err != nil {
			return nil, err
		}
		rawErrMsg, ok := badResponse["error"]
		if !ok {
			return nil, errors.New("failed to refresh token")
		}
		errMsg, ok := rawErrMsg.(string)
		if !ok {
			return nil, errors.New("failed to refresh token")
		}
		return nil, errors.New(errMsg)

	}

	var refreshResponse refreshResponse
	if err := decoder.Decode(&refreshResponse); err != nil {
		return refreshResponse, err
	}
	return refreshResponse, nil
}
