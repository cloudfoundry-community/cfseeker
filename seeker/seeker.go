package seeker

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cloudfoundry-community/cfseeker/config"
	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/cloudfoundry-community/gogobosh"
	"github.com/starkandwayne/goutils/log"
)

//Seeker has constructs and functions necessary to find app locations in Cloud
// Foundry
type Seeker struct {
	CF      *cfclient.Client
	bosh    *gogobosh.Client
	config  *config.Config
	vmcache *VMCache
}

//NewSeeker returns a NewSeeker with a client configured with the information
// in the given Config object
func NewSeeker(conf *config.Config) (ret *Seeker, err error) {
	ret = &Seeker{}
	ret.config = conf

	log.Debugf("Setting up CF Client")
	ret.CF, err = ret.getCFClientFromConfig()
	if err != nil {
		return nil, fmt.Errorf("Error connecting to Cloud Foundry API: %s", err.Error())
	}
	log.Debugf("Done setting up CF Client")

	if ret.BOSHConfigured() {
		log.Debugf("Setting up BOSH Client")
		ret.bosh, err = ret.getBOSHClientFromConfig()
		if err != nil {
			return nil, fmt.Errorf("Error connecting to BOSH API: %s", err.Error())
		}

		log.Debugf("Done setting up BOSH Client")
	} else {
		log.Debugf("Skipping BOSH Client setup")
	}

	ret.vmcache = newVMCache()
	return
}

// getClientFromConfig returns a cfclient.Client object that has been
// initialized with the settings that have been configured in the receiver
// Seeker object's Config
func (s *Seeker) getCFClientFromConfig() (client *cfclient.Client, err error) {
	return cfclient.NewClient(&cfclient.Config{
		ApiAddress:        s.config.CF.APIAddress,
		ClientID:          s.config.CF.ClientID,
		ClientSecret:      s.config.CF.ClientSecret,
		HttpClient:        &http.Client{Timeout: time.Second * time.Duration(s.config.HTTPTimeout)},
		SkipSslValidation: s.config.CF.SkipSSLValidation,
		UserAgent:         "Go-CF-client/1.1",
	})
}

// getBOSHClientFromConfig returns a gogobosh.Client object that has been
// initialized with the settings given in the config in the receiver Seeker
// object.
func (s *Seeker) getBOSHClientFromConfig() (client *gogobosh.Client, err error) {
	return gogobosh.NewClient(&gogobosh.Config{
		BOSHAddress:       s.config.BOSH.APIAddress,
		Username:          s.config.BOSH.Username,
		Password:          s.config.BOSH.Password,
		ClientID:          s.config.BOSH.ClientID,
		ClientSecret:      s.config.BOSH.ClientSecret,
		HttpClient:        &http.Client{Timeout: time.Second * time.Duration(s.config.HTTPTimeout)},
		SkipSslValidation: s.config.BOSH.SkipSSLValidation,
	})
}

// BOSHConfigured returns true if the attached configuration has all the keys
// required to attempt a connection to BOSH. False otherwise.
func (s *Seeker) BOSHConfigured() bool {
	conf := s.config.BOSH
	return conf.APIAddress != "" &&
		len(conf.Deployments) > 0 &&
		((conf.Username != "" && conf.Password != "") || (conf.ClientID != "" && conf.ClientSecret != "")) &&
		!conf.SkipBOSH
}
