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
	AppGUID   string         `yaml:"app_guid"`
	Instances []FindInstance `yaml:"instances"`
	Count     int            `yaml:"count"`
}

//FindInstance represents information about one instance of an app
type FindInstance struct {
	InstanceNumber int    `yaml:"number"`
	VMName         string `yaml:"vm_name"`
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
}

//Find determines the location of the app you requests
func Find(s *seeker.Seeker, in FindInput) (output FindOutput, err error) {
	log.Debugf("Beginning evaluation of find command")
	ret := FindOutput{}
	err = validateFindFlags(in)
	if err != nil {
		return
	}

	var instances []seeker.AppInstance

	if in.AppGUID != "" {
		log.Debugf("Finding IPs by GUID")
		ret.AppGUID = in.AppGUID
		instances, err = s.FindInstances(s.ByGUID(in.AppGUID))
	} else {
		log.Debugf("Finding IPs by Org, Space, and App Name")
		ret.AppGUID, err = /*Get App GUID*/ s.ByOrgSpaceAndName(in.OrgName, in.SpaceName, in.AppName)
		instances, err = s.FindInstances(ret.AppGUID, err)
	}

	if err != nil {
		err = fmt.Errorf("Error while getting VM IPs: %s", err.Error())
		return
	}

	log.Debugf("Got VM IPs")

	for i, instance := range instances {
		log.Debugf("Looking up VM with IP: %s", instance.Host)
		var vm *seeker.VMInfo
		vm, err = s.GetVMWithIP(instance.Host)

		if err != nil {
			err = fmt.Errorf("Error while translating VM name for IP `%s`: %s", instance.Host, err.Error())
			return
		}

		if vm == nil {
			//TODO: Don't error out, have alternate branch for no name resolution #resiliency
			err = fmt.Errorf("Could not find VM with given IP `%s`", instance.Host)
			return
		}

		log.Debugf("Got VM with IP: %s", instance.Host)

		thisInstance := FindInstance{
			InstanceNumber: i,
			VMName:         fmt.Sprintf("%s/%d", vm.Name, vm.Index),
			Host:           instance.Host,
			Port:           instance.Port,
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
