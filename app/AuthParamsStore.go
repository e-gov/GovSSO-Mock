package main

import (
	"math/rand"
	"strconv"
)

type authParamsStore struct {
	paramsByCode map[string]authParams
}

func newAuthParamsStore() *authParamsStore {
	return &authParamsStore{
		paramsByCode: make(map[string]authParams),
	}
}

func (this *authParamsStore) addParams(params authParams) string {
	code := strconv.Itoa(rand.Int())
	this.paramsByCode[code] = params
	return code
}

func (this *authParamsStore) getParams(code string) *authParams {
	if params, exists := this.paramsByCode[code]; exists {
		return &params
	}
	return nil
}
