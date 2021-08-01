//go:generate packer-sdc mapstructure-to-hcl2 -type AuthConfig

package common

import tss "github.com/thycotic/tss-sdk-go/server"

type AuthConfig struct {
	// Username of account with access to Thycotic Secret Server
	Username string `mapstructure:"username" require:"true"`
	// Password of account with the provided Username
	Password string `mapstructure:"password" require:"true"`
	// The Secret Server base URL. E.G. https://localhost/SecretServer
	ServerURL string `mapstructure:"server_url" require:"true"`
}

func (c *AuthConfig) CreateClient() (*tss.Server, error) {
	config := tss.Configuration{
		Credentials: tss.UserCredential{
			Username: c.Username,
			Password: c.Password,
		},
		ServerURL: c.ServerURL,

		// TLD and Tenant fields are included for completeness, but not currently handled.
		TLD:    "",
		Tenant: "",
	}

	client, err := tss.New(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
