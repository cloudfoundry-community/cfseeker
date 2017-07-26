package config

//Version of cfseeker gets set here
var Version = ""

func init() {
	if Version == "" {
		Version = "/shrug"
	}
}

//Config contains all the information needed for the seeker backend to operate
type Config struct {
	CF          CFConfig     `yaml:"cf"`
	BOSH        BOSHConfig   `yaml:"bosh"`
	Server      ServerConfig `yaml:"server"`
	HTTPTimeout int          `yaml:"http_timeout"`
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
	ClientID          string   `yaml:"client_id"`
	ClientSecret      string   `yaml:"client_secret"`
	SkipSSLValidation bool     `yaml:"skip_ssl_validation"`
	Deployments       []string `yaml:"deployments"`
	SkipBOSH          bool     `yaml:"skip_bosh"`
}

//ServerConfig has the info needed specifically for running in server mode
type ServerConfig struct {
	BasicAuth BasicAuthConfig `yaml:"basic_auth"`
	Port      int             `yaml:"port"`
	NoAuth    bool            `yaml:"no_auth"`
	CacheTTL  int             `yaml:"cache_ttl"` //in seconds
}

//BasicAuthConfig lets you set up basic auth for your API
type BasicAuthConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

//SkipBOSH sets the seeker config to not connect to BOSH
func (c *Config) SkipBOSH() {
	c.BOSH.SkipBOSH = true
}
