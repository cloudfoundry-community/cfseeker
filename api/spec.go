package api

const (
	//FindEndpoint is the URL endpoint corresponding to calling the Find command
	FindEndpoint = "/v1/apps"
	// MetaEndpoint is the URL endpoint corresponding to getting meta information
	// about this cfseeker server
	MetaEndpoint = "/v1/meta"
	// InvalidateBOSHEndpoint is the endpoint corresponding to manipulation of the
	// BOSH VM info cache
	InvalidateBOSHEndpoint = "/v1/cache/bosh"
	//WebEndpoint is the path to the web UI
	WebEndpoint = "/"
	//ConvertEndpoint is the path corresponding to the Convert API call
	ConvertEndpoint = "/v1/convert"
)
