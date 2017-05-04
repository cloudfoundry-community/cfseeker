package api

import (
	"net/http"

	"github.com/thomasmmitchell/cfseeker/config"
)

//MetaOutput gives meta information about this cfseeker server.
type MetaOutput struct {
	Version string `json:"version" yaml:"version"`
}

func metaHandler(w http.ResponseWriter, r *http.Request) {
	output := MetaOutput{
		Version: config.Version,
	}
	NewResponse(w).AttachContents(output).Write()
}
