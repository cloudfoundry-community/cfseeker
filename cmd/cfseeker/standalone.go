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
	}
	return
}

func findCommand(input interface{}) (interface{}, error) {
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

func serverCommand(input interface{}) (interface{}, error) {
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
