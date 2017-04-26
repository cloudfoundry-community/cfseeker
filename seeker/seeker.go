package seeker

import (
	"net/http"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/thomasmmitchell/cfseeker/config"
)

//Seeker has constructs and functions necessary to find app locations in Cloud
// Foundry
type Seeker struct {
	client *cfclient.Client
	config *config.Config
}

//NewSeeker returns a NewSeeker with a client configured with the information
// in the given Config object
func NewSeeker(conf *config.Config) (ret *Seeker, err error) {
	ret = &Seeker{}
	ret.config = conf
	err = ret.SetClientFromConfig()
	return
}

//SetClientFromConfig returns a cfclient.client object that has been initialized with
// the settings that have been configured in the config package.
func (s *Seeker) SetClientFromConfig() (err error) {
	s.client, err = cfclient.NewClient(&cfclient.Config{
		ApiAddress:        s.config.CF.APIAddress,
		ClientID:          s.config.CF.ClientID,
		ClientSecret:      s.config.CF.ClientSecret,
		HttpClient:        http.DefaultClient,
		SkipSslValidation: s.config.CF.SkipSslValidation,
		UserAgent:         "Go-CF-client/1.1",
	})
	return
}
