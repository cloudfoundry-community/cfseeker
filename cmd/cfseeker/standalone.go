package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cloudfoundry-community/cfseeker/api"
	"github.com/cloudfoundry-community/cfseeker/commands"
	"github.com/cloudfoundry-community/cfseeker/config"
	"github.com/cloudfoundry-community/cfseeker/seeker"
)

func getStandaloneFn(command string) (toRun commandFn, toInput interface{}) {
	switch command {
	case "find":
		toRun = findCommand
		toInput = commands.FindInput{
			AppGUID:   *appGUIDFind,
			OrgName:   *orgFind,
			SpaceName: *spaceFind,
			AppName:   *appNameFind,
		}
	case "server":
		toRun = serverCommand
		toInput = serverInput{conf: conf}
	case "invalidate":
		bailWith("Cannot run invalidate command without --target (-t) set")
	case "info", "meta":
		bailWith("Cannot run info command without --target (-t) set")
	case "convert guid":
		toRun = convertCommand
		toInput = commands.ConvertInput{
			GUID: *guidGUIDConv,
		}
	case "convert org":
		toRun = convertCommand
		toInput = commands.ConvertInput{
			OrgName: *orgNameOrgConv,
		}
	case "convert space":
		toRun = convertCommand
		toInput = commands.ConvertInput{
			OrgName:   *orgNameSpaceConv,
			SpaceName: *spaceNameSpaceConv,
		}
	case "convert app":
		toRun = convertCommand
		toInput = commands.ConvertInput{
			OrgName:   *orgNameAppConv,
			SpaceName: *spaceNameAppConv,
			AppName:   *appNameAppConv,
		}
	default:
		bailWith("Unrecognized command: %s", command)
	}
	return
}

func findCommand(input interface{}) (seeker.Output, error) {
	in := input.(commands.FindInput)
	s, err := seeker.NewSeeker(conf)
	if err != nil {
		return nil, err
	}
	return commands.Find(s, in)
}

type serverInput struct {
	conf *config.Config
}

func serverCommand(input interface{}) (seeker.Output, error) {
	var err error
	in := input.(serverInput)
	if *cfModeServer {
		portString := os.Getenv("PORT")
		in.conf.Server.Port, err = strconv.Atoi(portString)
		if err != nil {
			return nil, fmt.Errorf("PORT environment variable cannot be converted to int")
		}
	}
	err = api.Initialize(in.conf) //Never exits without an error
	return nil, err
}

func convertCommand(input interface{}) (seeker.Output, error) {
	var err error
	in := input.(commands.ConvertInput)
	conf.SkipBOSH()
	s, err := seeker.NewSeeker(conf)
	if err != nil {
		return nil, err
	}
	return commands.Convert(s, in)
}
