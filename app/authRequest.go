package main

type authRequest struct {
	ClientId     string `form:"client_id"`
	Nonce        string `form:"nonce"`
	RedirectUri  string `form:"redirect_uri"`
	ResponseType string `form:"response_type" binding:"valid_response_type"`
	Scope        string `form:"scope" binding:"valid_scope"`
	State        string `form:"state" binding:"valid_state"`
	UILocales    string `form:"ui_locales"`
	AcrValues    string `form:"acr_values" binding:"omitempty,oneof=low substantial high"`
}
