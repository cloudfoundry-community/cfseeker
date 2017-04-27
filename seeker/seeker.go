package seeker

import (
	"fmt"
	"net/http"
	"sync"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/cloudfoundry-community/gogobosh"
	"github.com/thomasmmitchell/cfseeker/config"
)

//Seeker has constructs and functions necessary to find app locations in Cloud
// Foundry
type Seeker struct {
	cf     *cfclient.Client
	bosh   *gogobosh.Client
	config *config.Config
	vmdata map[string]VMInfo
	//List of deployments in the BOSH Director who have their VMs currently cached
	cachedDeps []string
	lock       sync.Mutex
}

//NewSeeker returns a NewSeeker with a client configured with the information
// in the given Config object
func NewSeeker(conf *config.Config) (ret *Seeker, err error) {
	ret = &Seeker{}
	ret.config = conf

	ret.cf, err = ret.getCFClientFromConfig()
	if err != nil {
		return nil, fmt.Errorf("Error connecting to Cloud Foundry API: %s", err.Error())
	}

	ret.bosh, err = ret.getBOSHClientFromConfig()
	if err != nil {
		return nil, fmt.Errorf("Error connecting to BOSH API: %s", err.Error())
	}

	ret.vmdata = map[string]VMInfo{}
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
		HttpClient:        http.DefaultClient,
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
		HttpClient:        http.DefaultClient,
		SkipSslValidation: s.config.BOSH.SkipSSLValidation,
	})
}
