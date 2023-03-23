package main

type user struct {
	Subject     string `json:"sub"`
	GivenName   string `json:"given_name"`
	FamilyName  string `json:"family_name"`
	Birthdate   string `json:"birthdate"`
	Amr         string `json:"amr"`
	Acr         string `json:"acr"`
	PhoneNumber string `json:"phone_number"`
}
