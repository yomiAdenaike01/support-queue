package oauth

type IProvider interface {
	authorisationEndpoint() string
	getName() ProviderName
	setProperties(deps providerDeps)
	tokenExchange(code string) (ITokenResponse, error)
	refresh(refreshToken string) (iRefreshResponse, error)
}

type Provider struct {
	Name               ProviderName `json:"-"`
	Issuer             string       `json:"issuer"`
	AuthEndpoint       string       `json:"authorization_endpoint"`
	TokenEndpoint      string       `json:"token_endpoint"`
	RevocationEndpoint string       `json:"revocation_endpoint"`
	Scopes             string       `json:"-"`
	ClientId           string       `json:"-"`
	ClientSecret       string       `json:"-"`
	RedirectUrl        string       `json:"-"`
}

func (p Provider) getName() ProviderName {
	return p.Name
}
