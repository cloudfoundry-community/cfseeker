package api

import (
	"net/http"

	"github.com/thomasmmitchell/cfseeker/commands"
	"github.com/thomasmmitchell/cfseeker/seeker"
)

const (
	// FindAppGUIDKey is the HTTP query key for the App GUID to the Find API call.
	FindAppGUIDKey = "app_guid"
	// FindOrgNameKey is the HTTP query key for the Org Name to the Find API call.
	FindOrgNameKey = "org_name"
	// FindSpaceNameKey is the HTTP query key for the Space Name to the Find API
	// call.
	FindSpaceNameKey = "space_name"
	// FindAppNameKey is the HTTP query key for the App Name to the Find API call.
	FindAppNameKey = "app_name"
)

func findHandler(w http.ResponseWriter, r *http.Request, s *seeker.Seeker) {
	output, err := commands.Find(s, commands.FindInput{
		AppGUID:   r.FormValue(FindAppGUIDKey),
		OrgName:   r.FormValue(FindOrgNameKey),
		SpaceName: r.FormValue(FindSpaceNameKey),
		AppName:   r.FormValue(FindAppNameKey),
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
