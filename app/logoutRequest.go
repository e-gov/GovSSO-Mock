package main

type logoutRequest struct {
	IdTokenHint           string `form:"id_token_hint" binding:"required"`
	PostLogoutRedirectUri string `form:"post_logout_redirect_uri" binding:"required,url"`
	State                 string `form:"state"`
}
