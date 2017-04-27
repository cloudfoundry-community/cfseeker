package main

import (
	"fmt"
	"strings"

	"github.com/starkandwayne/goutils/log"
	"github.com/thomasmmitchell/cfseeker/seeker"
)

//FindOutput contains the return values from a call to Find()
type FindOutput struct {
	Instances []FindInstance `yaml:"instances"`
	Count     int            `yaml:"count"`
}

//FindInstance represents information about one instance of an app
type FindInstance struct {
	InstanceNumber int    `yaml:"number"`
	VMName         string `yaml:"vm_name"`
	Host           string `yaml:"host"`
}

func find(s *seeker.Seeker) (output interface{}, err error) {
	log.Debugf("Beginning evaluation of find command")
	ret := FindOutput{}
	err = validateFindFlags()
	if err != nil {
		return
	}

	var hosts []string

	if appGUIDFind != nil && *appGUIDFind != "" {
		log.Debugf("Finding IPs by GUID")
		hosts, err = s.FindIPs(s.ByGUID(*appGUIDFind))
	} else {
		log.Debugf("Finding IPs by Org, Space, and App Name")
		hosts, err = s.FindIPs(s.ByOrgSpaceAndName(*orgFind, *spaceFind, *appNameFind))
	}

	if err != nil {
		err = fmt.Errorf("Error while getting VM IPs: %s", err.Error())
		return
	}

	log.Debugf("Got VM IPs")

	for i, host := range hosts {
		log.Debugf("Looking up VM with IP: %s", host)
		var vm *seeker.VMInfo
		vm, err = s.GetVMWithIP(host)

		if err != nil {
			err = fmt.Errorf("Error while translating VM name for IP `%s`: %s", host, err.Error())
			return
		}

		log.Debugf("Got VM with IP: %s")

		thisInstance := FindInstance{
			InstanceNumber: i,
			VMName:         fmt.Sprintf("%s/%d", vm.Name, vm.Index),
			Host:           host,
		}
		ret.Instances = append(ret.Instances, thisInstance)
	}

	ret.Count = len(ret.Instances)

	output = ret
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
