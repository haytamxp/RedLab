package auth

import "strings"

type TokenProvider struct {
	token string
}

func NewTokenProvider(token string) *TokenProvider {
	return &TokenProvider{
		token: strings.TrimSpace(token),
	}
}

func (p *TokenProvider) Token() string {
	return p.token
}

func (p *TokenProvider) AuthorizationHeader() string {
	return "Bearer " + p.token
}