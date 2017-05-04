package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/starkandwayne/goutils/log"
	"github.com/thomasmmitchell/cfseeker/config"
	"github.com/thomasmmitchell/cfseeker/seeker"
)

//SeekerHandler is a modified version of http.HandlerFunc that receives a Seeker
// struct that can be used to perform operations
type SeekerHandler func(http.ResponseWriter, *http.Request, *seeker.Seeker)

var (
	configuration *config.Config
	defaultSeeker *seeker.Seeker
)

// Initialize reads in the given configuration struct and performs the steps
// necessary to prepare the server for launch
func Initialize(conf *config.Config) (err error) {
	log.Debugf("Beginning initialization of API server")
	configuration = conf
	err = validateServerConfig(conf.Server)
	if err != nil {
		return
	}

	var skipDefaultSeeker bool

	switch {
	case shouldBasicAuth(conf.Server):
		log.Debugf("Setting up basic auth")
		err = verifyBasicAuth(conf.Server)
		configuredAuth = basicAuth
	default:
		log.Debugf("Setting up no auth")
		err = verifyNopAuth(conf.Server)
		configuredAuth = nopAuth
	}
	//Check that the configured auth type is set up properly
	if err != nil {
		return fmt.Errorf("Error while configuring server auth: %s", err.Error())
	}

	//Eventually, I want to have a form of auth that just goes to the backend
	// CF UAA, in which case, we wouldn't need a default seeker
	if !skipDefaultSeeker {
		defaultSeeker, err = seeker.NewSeeker(conf)
		if err != nil {
			return fmt.Errorf("Error while creating seeker backend: %s", err.Error())
		}

		defaultSeeker.SetTTL(time.Duration(conf.Server.CacheTTL) * time.Second)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", auth(webHandler))
	router.HandleFunc(FindEndpoint, auth(findHandler)).Methods("GET")
	router.HandleFunc(MetaEndpoint, metaHandler).Methods("GET")

	router.NotFoundHandler = notFoundHandler{}

	log.Debugf("Listening on port %d", conf.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Server.Port), router)
	//If we're here, something is terrible
	return err
}

func validateServerConfig(conf config.ServerConfig) (err error) {
	if conf.Port > 65535 || conf.Port < 0 {
		err = fmt.Errorf("Port number %d is out of bounds", conf.Port)
	}
	return err
}
