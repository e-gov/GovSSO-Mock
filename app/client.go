package main

type client struct {
	ClientId               string   `json:"client_id"`
	BackchannelLogoutUri   string   `json:"backchannel_logout_uri"`
	RedirectUris           []string `json:"redirect_uris"`
	PostLogoutRedirectUris []string `json:"post_logout_redirect_uris"`
}

func (c client) isValidRedirectUri(redirectUri string) bool {
	return c.isRedirectUri(c.RedirectUris, redirectUri)
}

func (c client) isValidPostLogoutRedirectUri(redirectUri string) bool {
	return c.isRedirectUri(c.PostLogoutRedirectUris, redirectUri)
}

func (c client) isRedirectUri(validRedirectUris []string, redirectUri string) bool {
	if len(validRedirectUris) != 0 {
		for _, uri := range validRedirectUris {
			if uri == redirectUri {
				return true
			}
		}
	}
	return false
}
