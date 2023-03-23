package main

import "time"

type authParams struct {
	IdTokenClaims     map[string]interface{} `form:"id_token_claims"`
	LogoutTokenClaims map[string]interface{} `form:"logout_token_claims"`
	expires           time.Time
	tokenId           string
}

func (a *authParams) getSessionId() string {
	if a.IdTokenClaims["sid"] != nil {
		return a.IdTokenClaims["sid"].(string)
	} else {
		return ""
	}
}

func (a *authParams) isExpired() bool {
	return time.Now().After(a.expires)
}

func (a *authParams) isCustomLogoutToken() bool {
	return a.LogoutTokenClaims != nil
}
