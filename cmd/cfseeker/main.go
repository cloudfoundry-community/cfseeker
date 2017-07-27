package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry-community/cfseeker/config"
	"github.com/cloudfoundry-community/cfseeker/seeker"
	"github.com/starkandwayne/goutils/ansi"
	"github.com/starkandwayne/goutils/log"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
)

var (
	cmdLine = kingpin.New("cfseeker", "Do you know where your CF apps are?").Version(config.Version)
	//Global flags
	configPath   = cmdLine.Flag("config", "Path to a config file to load").Short('c').Default("./seekerconf.yml").Envar("SEEKERCONF").String()
	debugFlag    = cmdLine.Flag("debug", "Turn debug output on").Short('d').Bool()
	jsonFlag     = cmdLine.Flag("json", "Give output in JSON instead of YAML").Short('j').Bool()
	targetFlag   = cmdLine.Flag("target", "URL to target in CLI mode").Short('t').URL()
	usernameFlag = cmdLine.Flag("username", "Username for basic auth in CLI mode").Short('u').String()
	passwordFlag = cmdLine.Flag("password", "Password for basic auth in CLI mode. Will prompt if not given").Short('p').String()

	//FIND
	findCom     = cmdLine.Command("find", "Get the location of an app")
	orgFind     = findCom.Flag("org", "The organization where the app is pushed").Short('o').String()
	spaceFind   = findCom.Flag("space", "The space within the given org where the app is pushed").Short('s').String()
	appNameFind = findCom.Flag("app", "The name of the app to look up").Short('a').String()
	appGUIDFind = findCom.Flag("app-guid", "The GUID assigned to the app to look up").Short('g').String()

	//CONVERT
	convCom = cmdLine.Command("convert", "Convert from GUID to name")

	guidConvCom  = convCom.Command("guid", "Convert a GUID to the information it points to")
	guidGUIDConv = guidConvCom.Flag("guid", "GUID to get a name for").Short('g').Required().String()

	orgConvCom     = convCom.Command("org", "Convert an org name to its GUID")
	orgNameOrgConv = orgConvCom.Flag("org", "Name of the org").Short('o').Required().String()

	spaceConvCom       = convCom.Command("space", "Convert org and space names to its org and space GUIDs")
	orgNameSpaceConv   = spaceConvCom.Flag("org", "Name of the org").Short('o').Required().String()
	spaceNameSpaceConv = spaceConvCom.Flag("space", "Name of the space").Short('s').Required().String()

	appConvCom       = convCom.Command("app", "Convert org, space, and app names to their respective GUID information")
	orgNameAppConv   = appConvCom.Flag("org", "Name of the org").Short('o').Required().String()
	spaceNameAppConv = appConvCom.Flag("space", "Name of the space").Short('s').Required().String()
	appNameAppConv   = appConvCom.Flag("app", "Name of the app").Short('a').Required().String()

	//SERVER
	serverCom    = cmdLine.Command("server", "Run cfseeker in server mode")
	cfModeServer = serverCom.Flag("cf", "Override port in config to use PORT environment variable").Bool()

	//INVALIDATE
	invalidateCom = cmdLine.Command("invalidate", "Invalidate the BOSH cache on a cfseeker server")

	//INFO
	infoCom = cmdLine.Command("info", "Gives info about a running cfseeker server").Alias("meta")

	// //LIST
	// listCom = cmdLine.Command("list", "List all the apps on a given BOSH VM")
	// vmList  = listCom.Flag("vm", "The vm name to list instances for (<jobname>/<index>)").Required().String()
	conf *config.Config
)

type commandFn func(inputs interface{}) (seeker.Output, error)

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

	if targetIsSet() { //CLI mode
		validateTargetURI() //Bails if not valid
		toRun, toInput = getCLIFn(command)
	} else { //standalone or server mode
		toRun, toInput = getStandaloneFn(command)
	}

	log.Debugf("Dispatching to user command")
	cmdOut, err := toRun(toInput)
	if err != nil {
		bailWith(err.Error())
	}

	log.Debugf("Done with user command")

	var userOutput []byte

	if *jsonFlag {
		userOutput, err = json.Marshal(cmdOut)
	} else {
		userOutput, err = yaml.Marshal(cmdOut)
	}
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
	//Set defaults
	ret.Server.CacheTTL = 60 * 15 //15 Minutes
	ret.HTTPTimeout = 15          //15 seconds
	err = yaml.Unmarshal(configBytes, &ret)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing config YAML: %s", err.Error())
	}

	return &ret, nil
}

func targetIsSet() bool {
	return targetFlag != nil && *targetFlag != nil
}

func validateTargetURI() {
	if (*targetFlag).Scheme == "" {
		bailWith("Scheme (http / https) not given in target URI")
	}
	if (*targetFlag).Host == "" {
		bailWith("Host not given in target URI")
	}
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
