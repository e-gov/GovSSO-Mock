package main

import (
	"fmt"
)

type config struct {
	Host                      string `json:"host"`
	ServerPort                string `json:"serverPort"`
	BaseHref                  string `json:"baseHref"`
	TLSCertificate            string `json:"tlsCertificate"`
	TLSPrivateKey             string `json:"tlsPrivateKey"`
	IdTokenSignPrivateKeyPath string `json:"idTokenSignPrivateKeyPath"`
	IdTokenSignPublicKeyPath  string `json:"idTokenSignPublicKeyPath"`
	IdTokenSignKeyId          string `json:"idTokenSignKeyId"`
}

func (this config) HostUri() string {
	return fmt.Sprintf("https://%s:%s%s", this.Host, this.ServerPort, this.BaseHref)
}
