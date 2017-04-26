package config

var Conf Config

//Config contains all the information needed for the seeker backend to operate
type Config struct {
	CF CFConfig `yaml:"cf"`
}

//CFConfig contains location and authorization info about a target Cloud Foundry
type CFConfig struct {
	APIAddress        string `yaml:"api_address"`
	ClientID          string `yaml:"client_id"`
	ClientSecret      string `yaml:"client_secret"`
	SkipSslValidation bool   `yaml:"skip_ssl_validation"`
}
