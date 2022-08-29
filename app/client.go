package main

type client struct {
	ClientId             string `json:"client_id"`
	BackchannelLogoutUri string `json:"backchannel_logout_uri"`
}
