package commands

import (
	"fmt"
	"strings"

	"github.com/starkandwayne/goutils/log"
	"github.com/thomasmmitchell/cfseeker/seeker"
)

//FindInput contains the information required to perform the find command
// Either should be the org, space, and app names, or just the app GUID
type FindInput struct {
	AppGUID   string
	OrgName   string
	SpaceName string
	AppName   string
}

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

//Find determines the location of the app you requests
func Find(s *seeker.Seeker, in FindInput) (output FindOutput, err error) {
	log.Debugf("Beginning evaluation of find command")
	ret := FindOutput{}
	err = validateFindFlags(in)
	if err != nil {
		return
	}

	var hosts []string

	if in.AppGUID != "" {
		log.Debugf("Finding IPs by GUID")
		hosts, err = s.FindIPs(s.ByGUID(in.AppGUID))
	} else {
		log.Debugf("Finding IPs by Org, Space, and App Name")
		hosts, err = s.FindIPs(s.ByOrgSpaceAndName(in.OrgName, in.SpaceName, in.AppName))
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

		if vm == nil {
			//TODO: Don't error out, have alternate branch for no name resolution #resiliency
			err = fmt.Errorf("Could not find VM with given IP")
			return
		}

		log.Debugf("Got VM with IP: %s", host)

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

func validateFindFlags(in FindInput) error {
	//Check GUID flags
	if in.AppName != "" {
		return nil
	}

	var errorMessages []string
	//Otherwise, check org space name flags
	if in.OrgName == "" {
		errorMessages = append(errorMessages, "no org name specified")
	}
	if in.SpaceName == "" {
		errorMessages = append(errorMessages, "no space name specified")
	}
	if in.AppName == "" {
		errorMessages = append(errorMessages, "no app name specified")
	}

	if len(errorMessages) == 0 {
		return nil
	}

	return inputErrorf(strings.Join(errorMessages, "\n"))
}
