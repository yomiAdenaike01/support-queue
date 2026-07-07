package oauth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/yomiAdenaike01/support-queue/internals/config"
)

func newProvider(provider IProvider, deps providerDeps) IProvider {
	response, err := http.DefaultClient.Get(deps.wellKnownEndpoint)
	if err != nil {
		panic(err)
	}
	if err := json.NewDecoder(response.Body).Decode(provider); err != nil {
		panic(err)
	}
	provider.setProperties(deps)
	return provider

}

func getAuthConfig() authConfig {
	var configMap map[ProviderName]providerConfig
	if err := json.NewDecoder(bytes.NewReader(oauthConfig)).Decode(&configMap); err != nil {
		panic(err)
	}
	return configMap
}

func NewAuthProviders(config *config.Config, authIntegrationRepo AuthRepository) *Providers {
	authConfig := getAuthConfig()
	providerList := map[ProviderName]IProvider{
		GOOGLE: new(googleProvider),
	}
	register := map[ProviderName]IProvider{}

	providersChan := make(chan IProvider, len(providerList))
	var wg sync.WaitGroup
	for provider, providerStruct := range providerList {
		wg.Add(1)
		go func(provider ProviderName) {
			defer wg.Done()
			providersChan <- newProvider(providerStruct, authConfig.get(provider).getProviderDeps(config))
		}(provider)
	}
	go func() {
		for loadedProvider := range providersChan {
			register[loadedProvider.getName()] = loadedProvider
		}
	}()

	wg.Wait()

	return &Providers{
		register:             register,
		authIntegrationsRepo: authIntegrationRepo,
	}
}
