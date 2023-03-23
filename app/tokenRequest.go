package main

type tokenRequest struct {
	Code         string `form:"code"          binding:"required_without=RefreshToken"`
	RefreshToken string `form:"refresh_token" binding:"required_without=Code"`
	GrantType    string `form:"grant_type"    binding:"required,oneof=authorization_code refresh_token"`
}

func (t *tokenRequest) getId() string {
	if t.isRefreshTokenRequest() {
		return t.RefreshToken
	} else {
		return t.Code
	}
}

func (t *tokenRequest) isRefreshTokenRequest() bool {
	if len(t.RefreshToken) != 0 {
		return true
	} else {
		return false
	}
}
