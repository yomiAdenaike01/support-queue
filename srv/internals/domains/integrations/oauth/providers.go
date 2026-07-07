package oauth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yomiAdenaike01/support-queue/internals/config"
)

type ProviderName string

const (
	SLACK  ProviderName = "slack"
	GOOGLE ProviderName = "google"
)

type ITokenResponse interface {
	toSaveInput() saveTokenInput
}

type iRefreshResponse interface {
	getAccessToken() string
	getExpiresAt() time.Time
	toSaveInput() saveTokenInput
}

type providerConfig struct {
	WellKnownEndpoint string       `json:"wellKnownEndpoint"`
	Scopes            string       `json:"scopes"`
	ClientId          string       `json:"clientId"`
	ClientSecret      string       `json:"clientSecret"`
	Name              ProviderName `json:"name"`
}

type Callback struct {
	Code     string       `form:"code" binding:"required"`
	Provider ProviderName `form:"provider" binding:"required"`
}

type authConfig map[ProviderName]providerConfig

type Providers struct {
	register             map[ProviderName]IProvider
	authIntegrationsRepo AuthRepository
}

type providerDeps struct {
	name              ProviderName
	wellKnownEndpoint string
	scopes            string
	clientId          string
	clientSecret      string
	redirectUri       string
}

func (p *Providers) get(providerName ProviderName) IProvider {
	provider, ok := p.register[providerName]
	if !ok {
		log.Printf("Failed to find provider by name=%s", provider)
		return nil
	}
	return provider

}

func (p *Providers) GetAccessToken(providerName ProviderName) (string, error) {
	tokens, err := p.authIntegrationsRepo.GetTokens(context.TODO(), providerName)
	if err != nil {
		return "", err
	}
	accessToken, err := p.RefreshAccessToken(providerName, tokens)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (p *Providers) ExchangeAndSaveTokens(providerName ProviderName, code string) error {
	provider := p.get(providerName)
	tokens, err := provider.tokenExchange(code)
	if err != nil {
		return err
	}
	if err = p.authIntegrationsRepo.saveTokens(
		context.Background(),
		tokens.toSaveInput(),
	); err != nil {
		return err
	}

	return nil
}

func (p *Providers) GetAuthorisationEndpoint(providerName ProviderName) string {
	provider := p.get(providerName)
	if provider == nil {
		panic("failed to get provider")
	}
	return provider.authorisationEndpoint()
}

func (p *Providers) RefreshAccessToken(providerName ProviderName, tokens tokens) (string, error) {
	if !tokens.ShouldRefreshAccessToken() {
		return tokens.AccessToken, nil
	}
	provider := p.get(providerName)
	if provider == nil {
		panic("failed to get provider")
	}
	refresh, err := provider.refresh(tokens.RefreshToken)
	if err != nil {
		return "", err
	}
	if err = p.authIntegrationsRepo.saveTokens(context.TODO(), refresh.toSaveInput()); err != nil {
		return "", err
	}
	return refresh.getAccessToken(), nil
}

func (c providerConfig) getProviderDeps(config *config.Config) providerDeps {
	return providerDeps{
		name:              c.Name,
		redirectUri:       fmt.Sprintf("%s?provider=%s", config.GetEnvOrFail("OAUTH_CALLBACK_URL"), c.Name),
		wellKnownEndpoint: c.WellKnownEndpoint,
		scopes:            config.GetEnvOrFail(c.Scopes),
		clientId:          config.GetEnvOrFail(c.ClientId),
		clientSecret:      config.GetEnvOrFail(c.ClientSecret),
	}
}

func (a authConfig) get(providerName ProviderName) providerConfig {
	details, ok := a[providerName]
	if !ok {
		panic("failed to find provider")
	}
	return details
}
