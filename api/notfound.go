package api

import "net/http"

type notFoundHandler struct{}

func (notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	NewResponse(w).Code(404).Err("endpoint not found").Write()
}
