package infra

import (
	"errors"
	"github.com/dependabot/cli/internal/model"
)

// ConfigFilePath is the path to proxy config file.
const ConfigFilePath = "/config.json"

// Config is the structure of the proxy's config file
type Config struct {
	Credentials []model.Credential   `json:"all_credentials"`
	CA          CertificateAuthority `json:"ca"`
}

// CertificateAuthority includes the MITM CA certificate and private key
// Lite: RSA 3072 + secure constraints (v0.5.4-lite)
type CertificateAuthority struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

// BasicAuthCredentials represents credentials required for HTTP basic auth
type BasicAuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Lite: Validación para no levantar proxy sin CA
func (c Config) Validate() error {
	if c.CA.Cert == "" || c.CA.Key == "" {
		return errors.New("infra: CA cert/key empty - generate with cadetails.go")
	}
	// Permitimos creds vacías para runs públicos
	return nil
}
