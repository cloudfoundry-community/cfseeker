package api

import (
	"net/http"

	"github.com/cloudfoundry-community/cfseeker/commands"
	"github.com/cloudfoundry-community/cfseeker/seeker"
)

const (
	// ConvertGUIDKey is the HTTP query key for the GUID to search for in the
	// Convert API call
	ConvertGUIDKey = "guid"
	// ConvertOrgNameKey is the HTTP query key for the Org Name to search for in
	// the Convert API call
	ConvertOrgNameKey = "org_name"
	// ConvertSpaceNameKey is the HTTP query key for the Space Name to search for in the Convert API call
	// call.
	ConvertSpaceNameKey = "space_name"
	// ConvertAppNameKey is the HTTP query key for the App Name to search for in
	// the Convert API call.
	ConvertAppNameKey = "app_name"
)

func convertHandler(w http.ResponseWriter, r *http.Request, s *seeker.Seeker) {
	output, err := commands.Convert(s, commands.ConvertInput{
		GUID:      r.FormValue(ConvertGUIDKey),
		OrgName:   r.FormValue(ConvertOrgNameKey),
		SpaceName: r.FormValue(ConvertSpaceNameKey),
		AppName:   r.FormValue(ConvertAppNameKey),
	})

	if err != nil {
		if _, badRequest := err.(commands.InputError); badRequest {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		NewResponse(w).Err(err.Error()).Write()
		return
	}

	NewResponse(w).AttachContents(output).Write()
}
