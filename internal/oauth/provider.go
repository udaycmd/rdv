package oauth

import "golang.org/x/oauth2"

// Basic configuration for any Oauth2 provider
type BaseConfig struct {
	Name     string
	ClientId string
	Secret   string // https://developers.google.com/identity/protocols/oauth2/#installed
	Scopes   []string
	Ep       oauth2.Endpoint
}

// High level contract that should be implemented by every Oauth2 provider
type OauthProvider interface {
	GetInfo()
	GetCfg() *BaseConfig
}
