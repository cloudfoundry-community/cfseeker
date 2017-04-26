package main

import (
	"fmt"
	"strings"

	"github.com/thomasmmitchell/cfseeker/seeker"
)

func find(s *seeker.Seeker) (err error) {
	err = validateFindFlags()
	if err != nil {
		return
	}

	var host string

	if appGUIDFind != nil && *appGUIDFind != "" {
		host, err = s.FindIP(s.ByGUID(*appGUIDFind))
	} else {
		host, err = s.FindIP(s.ByOrgSpaceAndName(*orgFind, *spaceFind, *appNameFind))
	}

	//TODO: Print more betterer
	fmt.Println(host)
	//TODO: Go get vm name

	return
}

func validateFindFlags() error {
	//Check GUID flags
	if appGUIDFind != nil && *appGUIDFind != "" {
		return nil
	}

	var errorMessages []string
	//Otherwise, check org space name flags
	if orgFind == nil || *orgFind == "" {
		errorMessages = append(errorMessages, "No org name specified")
	}
	if spaceFind == nil || *spaceFind == "" {
		errorMessages = append(errorMessages, "No space name specified")
	}
	if appNameFind == nil || *appNameFind == "" {
		errorMessages = append(errorMessages, "No app name specified")
	}

	if len(errorMessages) == 0 {
		return nil
	}

	return fmt.Errorf(strings.Join(errorMessages, "\n"))
}
