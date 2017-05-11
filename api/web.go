package api

import (
	"net/http"

	"github.com/cloudfoundry-community/cfseeker/seeker"
)

func webHandler(w http.ResponseWriter, r *http.Request, _ *seeker.Seeker) {
	if r.URL.Path == "" || r.URL.Path == "/" {
		r.URL.Path = "/index.html"
	}

	if page, found := assets[r.URL.Path]; found {
		w.WriteHeader(http.StatusOK)
		w.Write(page)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("That page doesn't exist"))
	}
}
