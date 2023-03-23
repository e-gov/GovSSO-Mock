package main

import "time"

type authParamsStore struct {
	// authParams need to be stored and accessed by either `code` or `refresh_token`, depending on the specific use case.
	// As this is a mock application, we are storing them in a single map to keep things simpler.
	params map[string]authParams
	// in addition, authParams must be accessed by session id for some mocking use cases.
	paramsBySessionId map[string]*authParams
}

func newAuthParamsStore() *authParamsStore {
	return &authParamsStore{
		params:            make(map[string]authParams),
		paramsBySessionId: make(map[string]*authParams),
	}
}

func (this *authParamsStore) addParams(id string, params authParams) {
	params.tokenId = id
	params.expires = time.Now().Add(15 * time.Minute)
	this.params[id] = params
	sessionId := params.getSessionId()
	if len(sessionId) != 0 {
		this.paramsBySessionId[sessionId] = &params
	}
}

func (this *authParamsStore) getAndDeleteByTokenId(id string) *authParams {
	if params, exists := this.params[id]; exists {
		delete(this.params, id)
		delete(this.paramsBySessionId, params.getSessionId())
		return &params
	}
	return nil
}

func (this *authParamsStore) getAndDeleteBySessionId(id string) *authParams {
	if params, exists := this.paramsBySessionId[id]; exists {
		delete(this.params, params.tokenId)
		delete(this.paramsBySessionId, params.getSessionId())
		return params
	}
	return nil
}
