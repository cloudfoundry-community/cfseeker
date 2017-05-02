package api

import (
	"net/http"

	"github.com/thomasmmitchell/cfseeker/seeker"
)

func webHandler(w http.ResponseWriter, r *http.Request, _ *seeker.Seeker) {
	//TODO: Write this
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to CFSeeker. This will be a useable webpage in a future release."))
}
