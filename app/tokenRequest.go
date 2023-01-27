package main

type tokenRequest struct {
	Code      string `form:"code" binding:"required"`
	GrantType string `form:"grant_type" binding:"required,eq=authorization_code"`
}
