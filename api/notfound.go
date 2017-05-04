package api

import "net/http"

type notFoundHandler struct{}

func (notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write(NewResponse().Err("endpoint not found").Bytes())
}
