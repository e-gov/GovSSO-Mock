package main

type authParamsStore struct {
	// authParams need to be stored and accessed by either `code` or `refresh_token`, depending on the specific use case.
	// As this is a mock application, we are storing them in a single map to keep things simpler.
	params map[string]authParams
}

func newAuthParamsStore() *authParamsStore {
	return &authParamsStore{
		params: make(map[string]authParams),
	}
}

func (this *authParamsStore) addParams(id string, params authParams) {
	params.expires = this.getNewTokenExpirationTime()
	this.params[id] = params
}

func (this *authParamsStore) getAndDeleteParams(id string) *authParams {
	if params, exists := this.params[id]; exists {
		delete(this.params, id)
		return &params
	}
	return nil
}
