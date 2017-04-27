package config

//Config contains all the information needed for the seeker backend to operate
type Config struct {
	CF   CFConfig   `yaml:"cf"`
	BOSH BOSHConfig `yaml:"bosh"`
}

//CFConfig contains location and authorization info about a target Cloud Foundry
type CFConfig struct {
	APIAddress        string `yaml:"api_address"`
	ClientID          string `yaml:"client_id"`
	ClientSecret      string `yaml:"client_secret"`
	SkipSSLValidation bool   `yaml:"skip_ssl_validation"`
}

//BOSHConfig contains location, auth, and tracking info for your BOSH.
type BOSHConfig struct {
	APIAddress        string   `yaml:"api_address"`
	Username          string   `yaml:"username"`
	Password          string   `yaml:"password"`
	SkipSSLValidation bool     `yaml:"skip_ssl_validation"`
	Deployments       []string `yaml:"deployments"`
}
