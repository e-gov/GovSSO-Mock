package main

type authParams struct {
	acr        string
	amr        string
	birthdate  string
	clientId   []string
	givenName  string
	familyName string
	nonce      string
	scope      string
	state      string
	subject    string
	sessionId  *string
	phone      string
}
