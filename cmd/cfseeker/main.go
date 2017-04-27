package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/starkandwayne/goutils/ansi"
	"github.com/thomasmmitchell/cfseeker/config"
	"github.com/thomasmmitchell/cfseeker/seeker"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
)

var (
	cmdLine = kingpin.New("cf-seeker", "Do you know where your CF apps are?").Version("/shrug")
	//Global flags
	configPath = cmdLine.Flag("config", "Path to a config file to load").Short('c').Default("./seekerconf.yml").Envar("SEEKERCONF").String()

	//FIND
	findCom     = cmdLine.Command("find", "Get the location of an app")
	orgFind     = findCom.Flag("org", "The organization where the app is pushed").Short('o').String()
	spaceFind   = findCom.Flag("space", "The space within the given org where the app is pushed").Short('s').String()
	appNameFind = findCom.Flag("app", "The name of the app to look up").Short('a').String()
	appGUIDFind = findCom.Flag("appGUID", "The GUID assigned to the app to look up").Short('g').String()

	// //LIST
	// listCom = cmdLine.Command("list", "List all the apps on a given BOSH VM")
	// vmList  = listCom.Flag("vm", "The vm name to list instances for (<jobname>/<index>)").Required().String()
)

func main() {
	command := kingpin.MustParse(cmdLine.Parse(os.Args[1:]))
	cmdLine.HelpFlag.Short('h')
	cmdLine.VersionFlag.Short('v')
	conf, err := initializeConfig()
	if err != nil {
		bailWith(err.Error())
	}

	s, err := seeker.NewSeeker(conf)
	if err != nil {
		bailWith(err.Error())
	}

	switch command {
	case "find":
		err = find(s)
		// case "list":
	}

	if err != nil {
		bailWith(err.Error())
	}
}

func initializeConfig() (*config.Config, error) {
	ansi.Fprintf(os.Stderr, "@G{Using config path: %s}\n", *configPath)

	configFile, err := os.Open(*configPath)
	if err != nil {
		return nil, fmt.Errorf("Error opening config file: %s", err.Error())
	}

	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, fmt.Errorf("Error when reading config: %s", err.Error())
	}

	var ret config.Config
	err = yaml.Unmarshal(configBytes, &ret)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing config YAML: %s", err.Error())
	}

	return &ret, nil
}

func bailWith(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, ansi.Sprintf("@R{%s}\n", message), args...)
	os.Exit(1)
}
