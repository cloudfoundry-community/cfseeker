package commands

import (
	"encoding/json"
	"fmt"

	"github.com/cloudfoundry-community/cfseeker/seeker"
	"github.com/starkandwayne/goutils/log"
)

//ConvertInput is all of the information required for a call to Convert.
// If converting from Name to GUID, you must provide all names down to the level
// you are requesting information for. i.e:
// For Org name to GUID, you only need to provide org name
// For Space name to GUID, you need to provide org name and space name
// For App name to GUID, provide all three names.
//For converting from GUID to name, you only need to give the GUID because
// GUIDs are... globally unique (go figure). The function will search for what
// type of thing it belongs to.
type ConvertInput struct {
	//GUID is the GUID of the object to get the name of
	GUID string
	//OrgName is the name of the org to get the GUID of, or the org in which the
	// space or app you're fetching information for resides
	OrgName string
	//SpaceName is the name of the space to get the GUID of, or the space in which
	// the app you're fetching information for resides. If this is provided,
	// OrgName must also be provided.
	SpaceName string
	//AppName is the name of the app to get the GUID of. If this is provided,
	// SpaceName and OrgName must also be provided.
	AppName string
}

func (c ConvertInput) shouldConvGUID() bool {
	return c.GUID != "" && c.OrgName == "" && c.SpaceName == "" && c.AppName == ""
}

func (c ConvertInput) shouldConvOrg() bool {
	return c.GUID == "" && c.OrgName != "" && c.SpaceName == "" && c.AppName == ""
}

func (c ConvertInput) shouldConvSpace() bool {
	return c.GUID == "" && c.OrgName != "" && c.SpaceName != "" && c.AppName == ""
}

func (c ConvertInput) shouldConvApp() bool {
	return c.GUID == "" && c.OrgName != "" && c.SpaceName != "" && c.AppName != ""
}

//Validate returns true if the values given in the ConvertInput struct are a
// combination valid to make a request.
func (c ConvertInput) Validate() (err error) {
	if !(c.shouldConvGUID() || c.shouldConvOrg() || c.shouldConvSpace() || c.shouldConvApp()) {
		err = inputErrorf("Invalid combination of convert input arguments")
	}
	return
}

//ConvertOutput is a struct representing the information returned from a call
// to convert.
type ConvertOutput struct {
	//OrgGUID is the GUID of the org requested
	OrgGUID string `yaml:"org_guid" json:"org_guid"`
	//SpaceGUID is the GUID of the space requested, given back if space or app was
	// requested.
	SpaceGUID string `yaml:"space_guid,omitempty" json:"space_guid,omitempty"`
	//AppGUID is the GUID of the app requested, given back only if app info
	// was requested
	AppGUID string `yaml:"app_guid,omitempty" json:"app_guid,omitempty"`
	//OrgName is the name of the org requested
	OrgName string `yaml:"org_name" json:"org_name"`
	//SpaceName is the name of the space requested, given back only if space or
	// app info was requested
	SpaceName string `yaml:"space_name,omitempty" json:"space_name,omitempty"`
	//AppName is the name of the app requested, given back only if app info was
	// requested
	AppName string `yaml:"app_name,omitempty" json:"app_name,omitempty"`
	//Type of resource returned. One of `app`, `space`, or `org`
	Type string `yaml:"type" json:"type"`
}

//ReceiveJSON allows this to implement SeekerOutput
func (c *ConvertOutput) ReceiveJSON(j []byte) (err error) {
	err = json.Unmarshal(j, c)
	return
}

var (
	//ConvertTypeOrg indicates that the output returned represents an org
	ConvertTypeOrg = "org"
	//ConvertTypeSpace indicates that the output returned represents a space
	ConvertTypeSpace = "space"
	//ConvertTypeApp indicates that the output returned represents an app
	ConvertTypeApp = "app"
)

//Convert takes the input of the names of spaces, orgs, and/or apps and gives you
// the GUIDs associated with each. Alternatively, give it a GUID, and it gives you
// the names of the org/space/app which the GUID represents, and the names/GUIDs
// of the resources above it with their GUIDs and names
func Convert(s *seeker.Seeker, in ConvertInput) (out *ConvertOutput, err error) {
	err = in.Validate()
	if err != nil {
		return
	}

	out = &ConvertOutput{}

	switch {
	case in.shouldConvGUID():
		out, err = convGUID(s, in)
	case in.shouldConvOrg():
		out, err = convOrg(s, in)
	case in.shouldConvSpace():
		out, err = convSpace(s, in)
	case in.shouldConvApp():
		out, err = convApp(s, in)
	default:
		panic("Validated input did not correspond to a code path in Convert")
	}

	return
}

func convGUID(s *seeker.Seeker, in ConvertInput) (out *ConvertOutput, err error) {
	log.Debugf("Beginning conversion lookup by GUID")
	out = &ConvertOutput{}

	if out, err = convAppByGUID(s, in); err == nil {
	} else if out, err = convSpaceByGUID(s, in); err == nil {
	} else if out, err = convOrgByGUID(s, in); err == nil {
	} else {
		log.Debugf("All conversion lookups lookups failed")
		err = fmt.Errorf("Could not look up GUID: %s (does the GUID exist?)", in.GUID)
	}

	return
}

func convOrgByGUID(s *seeker.Seeker, in ConvertInput) (out *ConvertOutput, err error) {
	log.Debugf("Getting org by GUID")
	out = &ConvertOutput{}
	org, err := s.CF.GetOrgByGuid(in.GUID)
	if err != nil {
		err = fmt.Errorf("Error getting CF Org by GUID: %s", err.Error())
		return
	}

	out.OrgGUID = in.GUID
	out.OrgName = org.Name
	out.Type = ConvertTypeOrg

	log.Debugf("Successful org lookup by GUID")
	return
}

func convSpaceByGUID(s *seeker.Seeker, in ConvertInput) (out *ConvertOutput, err error) {
	log.Debugf("Getting space by GUID")
	out = &ConvertOutput{}
	space, err := s.CF.GetSpaceByGuid(in.GUID)
	if err != nil {
		err = fmt.Errorf("Error getting CF Space by GUID: %s", err.Error())
		return
	}

	out.SpaceGUID = in.GUID
	out.SpaceName = space.Name

	log.Debugf("Getting org associated with space with GUID (%s)", in.GUID)
	org, err := space.Org()
	if err != nil {
		err = fmt.Errorf("Error getting CF Org associated with space with GUID: %s", in.GUID)
		return
	}

	out.OrgGUID = org.Guid
	out.OrgName = org.Name
	out.Type = ConvertTypeSpace

	log.Debugf("Successful space lookup by GUID")
	return
}

func convAppByGUID(s *seeker.Seeker, in ConvertInput) (out *ConvertOutput, err error) {
	log.Debugf("Getting app by GUID")
	out = &ConvertOutput{}
	app, err := s.CF.GetAppByGuid(in.GUID)
	if err != nil {
		err = fmt.Errorf("Error getting CF App by GUID: %s", err.Error())
		return
	}

	out.AppGUID = in.GUID
	out.AppName = app.Name

	log.Debugf("Getting space associated with app with GUID (%s)", in.GUID)
	space, err := app.Space()
	if err != nil {
		err = fmt.Errorf("Error getting CF Space associated with app with GUID: %s", in.GUID)
		return
	}

	out.SpaceGUID = space.Guid
	out.SpaceName = space.Name

	log.Debugf("Getting org associated with space with GUID (%s)", space.Guid)
	org, err := space.Org()
	if err != nil {
		err = fmt.Errorf("Error getting CF Org associated with space with GUID: %s", space.Guid)
		return
	}

	out.OrgGUID = org.Guid
	out.OrgName = org.Name
	out.Type = ConvertTypeApp

	log.Debugf("Successful app lookup by GUID")
	return
}

func convOrg(s *seeker.Seeker, in ConvertInput) (out *ConvertOutput, err error) {
	log.Debugf("Beginning org conversion lookup")
	out = &ConvertOutput{}
	log.Debugf("Getting org by name (%s)", in.OrgName)
	org, err := s.CF.GetOrgByName(in.OrgName)
	if err != nil {
		err = fmt.Errorf("Error getting CF Org information: %s", err.Error())
		return
	}

	out.OrgName = in.OrgName
	out.OrgGUID = org.Guid
	out.Type = ConvertTypeOrg
	return
}

func convSpace(s *seeker.Seeker, in ConvertInput) (out *ConvertOutput, err error) {
	log.Debugf("Beginning space conversion lookup")
	out = &ConvertOutput{}
	//Need the GUID for the org
	out, err = convOrg(s, in)
	if err != nil {
		return
	}

	log.Debugf("Getting space by name (%s), and org GUID (%s)", in.SpaceName, out.OrgGUID)
	space, err := s.CF.GetSpaceByName(in.SpaceName, out.OrgGUID)
	if err != nil {
		err = fmt.Errorf("Error getting CF Space information: %s", err.Error())
		return
	}

	out.SpaceName = in.SpaceName
	out.SpaceGUID = space.Guid
	out.Type = ConvertTypeSpace
	return
}

func convApp(s *seeker.Seeker, in ConvertInput) (out *ConvertOutput, err error) {
	log.Debugf("Beginning app conversion lookup")
	out = &ConvertOutput{}
	//Need the GUID for the org and space. convSpace does org lookup for us.
	out, err = convSpace(s, in)
	if err != nil {
		return
	}

	log.Debugf("Getting app by name (%s), space GUID (%s), and org GUID (%s)", in.AppName, out.SpaceGUID, out.OrgGUID)
	app, err := s.CF.AppByName(in.AppName, out.SpaceGUID, out.OrgGUID)
	if err != nil {
		err = fmt.Errorf("Error getting CF App information: %s", err.Error())
		return
	}

	out.AppName = in.AppName
	out.AppGUID = app.Guid
	out.Type = ConvertTypeApp
	return
}
