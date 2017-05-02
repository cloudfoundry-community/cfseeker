package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/starkandwayne/goutils/ansi"
	"github.com/starkandwayne/goutils/log"
	"github.com/thomasmmitchell/cfseeker/commands"
	"github.com/thomasmmitchell/cfseeker/config"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
)

var (
	cmdLine = kingpin.New("cf-seeker", "Do you know where your CF apps are?").Version("/shrug")
	//Global flags
	configPath = cmdLine.Flag("config", "Path to a config file to load").Short('c').Default("./seekerconf.yml").Envar("SEEKERCONF").String()
	debugFlag  = cmdLine.Flag("debug", "Turn debug output on").Short('d').Bool()

	//FIND
	findCom     = cmdLine.Command("find", "Get the location of an app")
	orgFind     = findCom.Flag("org", "The organization where the app is pushed").Short('o').String()
	spaceFind   = findCom.Flag("space", "The space within the given org where the app is pushed").Short('s').String()
	appNameFind = findCom.Flag("app", "The name of the app to look up").Short('a').String()
	appGUIDFind = findCom.Flag("app-guid", "The GUID assigned to the app to look up").Short('g').String()

	//SERVER
	serverCom    = cmdLine.Command("server", "Run cfseeker in server mode")
	cfModeServer = cmdLine.Flag("cf", "Override port in config to use PORT environment variable").Bool()
	// //LIST
	// listCom = cmdLine.Command("list", "List all the apps on a given BOSH VM")
	// vmList  = listCom.Flag("vm", "The vm name to list instances for (<jobname>/<index>)").Required().String()
	conf *config.Config
)

func main() {
	cmdLine.HelpFlag.Short('h')
	cmdLine.VersionFlag.Short('v')

	command := kingpin.MustParse(cmdLine.Parse(os.Args[1:]))

	var err error
	conf, err = initializeConfig()
	if err != nil {
		bailWith(err.Error())
	}

	setupLogging()

	var toRun commandFn
	var toInput interface{}

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
	}

	log.Debugf("Dispatching to user command")
	cmdOut, err := toRun(toInput)
	if err != nil {
		bailWith(err.Error())
	}

	log.Debugf("Done with user command")

	userOutput, err := yaml.Marshal(cmdOut)
	if err != nil {
		bailWith("Could not marshal output into YAML")
	}

	fmt.Println(string(userOutput))
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

func setupLogging() {
	logLevel := "emerg"
	if *debugFlag {
		logLevel = "debug"
	}

	log.SetupLogging(log.LogConfig{
		Type:  "console",
		Level: logLevel,
	})
}

func bailWith(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, ansi.Sprintf("@R{%s}\n", message), args...)
	os.Exit(1)
}
