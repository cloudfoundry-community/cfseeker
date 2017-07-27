package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudfoundry-community/cfseeker/seeker"
	"github.com/starkandwayne/goutils/log"
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
	AppGUID   string         `yaml:"guid" json:"guid"`
	AppName   string         `yaml:"name" json:"name"`
	Instances []FindInstance `yaml:"instances" json:"instances"`
	Count     int            `yaml:"count" json:"count"`
}

//ReceiveJSON makes FindOutput an implementation of SeekerOutput
func (f *FindOutput) ReceiveJSON(j []byte) (err error) {
	err = json.Unmarshal(j, f)
	return
}

//FindInstance represents information about one instance of an app
type FindInstance struct {
	InstanceNumber int    `yaml:"number" json:"number"`
	VMName         string `yaml:"vm_name,omitempty" json:"vm_name,omitempty"`
	Deployment     string `yaml:"deployment,omitempty" json:"deployment,omitempty"`
	Host           string `yaml:"host" json:"host"`
	Port           int    `yaml:"port" json:"port"`
}

//Find determines the location of the app you requests
func Find(s *seeker.Seeker, in FindInput) (output *FindOutput, err error) {
	log.Debugf("Beginning evaluation of find command")
	ret := FindOutput{}
	err = validateFindFlags(in)
	if err != nil {
		return
	}

	var meta *seeker.AppMeta
	var instances []seeker.AppInstance

	if in.AppGUID != "" {
		log.Debugf("Finding IPs by GUID")
		ret.AppGUID = in.AppGUID
		meta, instances, err = s.FindInstances(s.ByGUID(in.AppGUID))
	} else {
		log.Debugf("Finding IPs by Org, Space, and App Name")
		meta, instances, err = s.FindInstances(s.ByOrgSpaceAndName(in.OrgName, in.SpaceName, in.AppName))
	}

	if err != nil {
		err = fmt.Errorf("Error while getting VM IPs: %s", err.Error())
		return
	}

	log.Debugf("Got VM IPs")

	ret.AppGUID = meta.GUID
	ret.AppName = meta.Name

	for i, instance := range instances {
		ret.Instances = append(ret.Instances, FindInstance{
			InstanceNumber: i,
			Host:           instance.Host,
			Port:           instance.Port,
		})
	}

	if s.BOSHConfigured() {
		lookupAndAssignBOSHInfo(ret.Instances, s)
	}

	ret.Count = len(ret.Instances)

	output = &ret
	return
}

func lookupAndAssignBOSHInfo(instances []FindInstance, s *seeker.Seeker) (err error) {
	for i, instance := range instances {
		log.Debugf("Looking up VM with IP: %s", instance.Host)
		var vm *seeker.VMInfo
		vm, err = s.GetVMWithIP(instance.Host)

		if err != nil {
			err = fmt.Errorf("Error while translating VM name for IP `%s`: %s", instance.Host, err.Error())
			return
		}

		if vm == nil {
			err = fmt.Errorf("Could not find VM with given IP `%s`", instance.Host)
			return
		}

		log.Debugf("Got VM with IP: %s", instance.Host)

		instances[i].Deployment = vm.DeploymentName
		instances[i].VMName = fmt.Sprintf("%s/%d", vm.JobName, vm.Index)
	}
	return
}

func validateFindFlags(in FindInput) error {
	//Check GUID flags
	if in.AppGUID != "" {
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
