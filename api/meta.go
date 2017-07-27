package api

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-community/cfseeker/config"
)

//MetaOutput gives meta information about this cfseeker server.
type MetaOutput struct {
	Version string `json:"version" yaml:"version"`
}

//ReceiveJSON makes MetaOutput an implementation of SeekerOutput
func (m *MetaOutput) ReceiveJSON(j []byte) (err error) {
	err = json.Unmarshal(j, m)
	return
}

func metaHandler(w http.ResponseWriter, r *http.Request) {
	output := &MetaOutput{
		Version: config.Version,
	}
	NewResponse(w).AttachContents(output).Write()
}
