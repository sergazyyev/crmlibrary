package ocrmconfigs

import (
	"context"
	"encoding/base64"
)

type BasicAuthenticateConf struct {
	Username string `toml:"basic_auth_username"`
	Password string `toml:"basic_auth_password"`
}

func (a *BasicAuthenticateConf) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	auth := a.Username + ":" + a.Password
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	return map[string]string{
		"authorization": "Basic " + enc,
	}, nil
}

func (a *BasicAuthenticateConf) RequireTransportSecurity() bool {
	return false
}

type UserCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *UserCredential) GetBasicAuthHeader() map[string]string {
	auth := a.Username + ":" + a.Password
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	return map[string]string{
		"authorization": "Basic " + enc,
	}
}

func (a *UserCredential) GetBasicAuthToken() string {
	auth := a.Username + ":" + a.Password
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	return enc
}
