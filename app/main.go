package main

import (
	"GOVSSO-Mock/app/idtoken"
	"GOVSSO-Mock/app/json"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/go-jose/go-jose/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

func main() {
	log.Info().Msg("Starting GovSSO-Mock")

	if config, err := json.LoadFile[config]("config/config.json"); err != nil {
		log.Fatal().Err(err).Msg("Failed to load app config")

	} else if idTokenSigningKey, err := loadIdTokenSigningKey(config); err != nil {
		log.Fatal().Err(err).Msg("Failed to load identity token signing key")

	} else if predefinedUsers, err := json.LoadFile[[]user]("config/users.json"); err != nil {
		log.Fatal().Err(err).Msg("Failed to load predefined users from file")

	} else if predefinedClients, err := json.LoadFile[[]client]("config/clients.json"); err != nil {
		log.Fatal().Err(err).Msg("Failed to load predefined clients from file")

	} else {
		handler := routeHandler{
			config:            config,
			idTokenSigningKey: idTokenSigningKey,
			predefinedUsers:   predefinedUsers,
			predefinedClients: predefinedClients,
			authParamsStore:   newAuthParamsStore(),
			idTokenService:    idtoken.NewService(idTokenSigningKey, config.HostUri(), config.IdTokenSignKeyId),
			httpClient:        tlsChecksDisabledHttpClient(),
		}
		if err := handler.init(); err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize route handler")
		}
	}
}

func loadIdTokenSigningKey(config config) (*jose.SigningKey, error) {
	privateKeyBytes, err := os.ReadFile(config.IdTokenSignPrivateKeyPath)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode(privateKeyBytes)
	if pemBlock == nil {
		return nil, errors.New("Unable to parse private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return &jose.SigningKey{Algorithm: jose.RS256, Key: privateKey}, nil
}

func tlsChecksDisabledHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
