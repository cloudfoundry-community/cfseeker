package api

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/starkandwayne/goutils/log"
	"github.com/thomasmmitchell/cfseeker/config"
)

type authorizer func(SeekerHandler) http.HandlerFunc

var (
	configuredAuth authorizer
)

func auth(h SeekerHandler) http.HandlerFunc {
	//Serve wants a HandlerFunc, so we give it an altered one that calls our auth
	// before calling the actual handler. Set the content-type to application/json
	// before deferring to the configured auth type, just so that doesn't need to
	// be done in all the other handlers
	return func(w http.ResponseWriter, request *http.Request) {
		configuredAuth(h)(w, request) //Form the handler with our auth, then call that handler
	}
}

//No auth - this is just a passthrough to the given HandlerFunc
func nopAuth(h SeekerHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		h(w, request, defaultSeeker)
	}
}

func verifyNopAuth(conf config.ServerConfig) (err error) {
	if !conf.NoAuth {
		noAuthField, ok := reflect.TypeOf(conf).FieldByName("NoAuth")
		if !ok {
			panic("Expected NoAuth field in server config")
		}
		err = fmt.Errorf("Cowardly refusing to set No Auth because %s is not set in the server section of the configuration file", noAuthField.Tag.Get("yaml"))

	}
	return
}

func basicAuth(h SeekerHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		//The easy part of basic auth
		isAuthorized := func(username, password string) bool {
			return username == configuration.Server.BasicAuth.Username &&
				password == configuration.Server.BasicAuth.Password
		}

		//Get basic auth if its there
		reqUser, reqPass, isBasicAuth := request.BasicAuth()
		if !isBasicAuth {
			errorMessage := "Authorization Failed: No Basic Auth Header"
			log.Infof("basicAuth: %s", errorMessage)
			w.Header().Set("WWW-Authenticate", "Basic realm=\"Portcullis API\"")

			w.WriteHeader(http.StatusUnauthorized)
			NewResponse(w).Err(errorMessage).Write()
			return
		}

		//Check the provided auth creds to see if they are what we should allow
		if !isAuthorized(reqUser, reqPass) {
			errorMessage := "Authorization Failed: Incorrect credentials"
			log.Warnf("basicAuth: %s", errorMessage)

			w.WriteHeader(http.StatusUnauthorized)
			NewResponse(w).Err(errorMessage).Write()
			return
		}
		h(w, request, defaultSeeker)
	}
}

func shouldBasicAuth(conf config.ServerConfig) bool {
	return conf.BasicAuth.Username != "" || conf.BasicAuth.Password != ""
}

func verifyBasicAuth(conf config.ServerConfig) (err error) {
	if conf.BasicAuth.Username == "" {
		err = fmt.Errorf("No basic auth username given")
	}
	if conf.BasicAuth.Password == "" {
		err = fmt.Errorf("No basic auth password given")
	}
	return
}
