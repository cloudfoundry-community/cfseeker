package seeker

import (
	"fmt"
	"net"

	"github.com/starkandwayne/goutils/log"
)

// AppMeta has information about the pushed application itself, universal to
// all its instances
type AppMeta struct {
	Name string
	GUID string
}

//AppInstance has information from the CF API about an application
type AppInstance struct {
	Host string
	Port int
}

// FindInstances takes the App GUID given and queries the CF API to get the IP
// addresses of the VMs and listening ports on which the instances of the
// application are located and returns those. If inputErr is non-nil, the
// function will bail early with the given error.
func (s *Seeker) FindInstances(guid string, inputErr error) (meta *AppMeta, inst []AppInstance, err error) {
	if inputErr != nil {
		return nil, nil, inputErr
	}

	meta = &AppMeta{}

	log.Debugf("Getting application stats for app with GUID %s from CF API", guid)
	statsMap, err := s.CF.GetAppStats(guid)
	if err != nil {
		err = fmt.Errorf("Error when getting stats for app with GUID `%s` (Is the app running?)", guid)
		return
	}
	if len(statsMap) == 0 {
		err = fmt.Errorf("No stats found for app with GUID `%s`", guid)
		return
	}

	meta.GUID = guid

	for _, stats := range statsMap {
		stats.Stats.Host, err = canonizeIP(stats.Stats.Host)
		if err != nil {
			return
		}
		inst = append(inst, AppInstance{Host: stats.Stats.Host, Port: stats.Stats.Port})

		meta.Name = stats.Stats.Name
	}

	return
}

// ByOrgSpaceAndName checks that the given variables are set, erroring if any
// of them are not, and then looks up the GUID of the app using the CF API.
func (s *Seeker) ByOrgSpaceAndName(org, space, app string) (retGUID string, err error) {
	return s.getAppGUID(org, space, app)
}

// ByGUID checks that the given appGUID is set, returning an error if not, and
// then passes through the given GUID, dereferenced.
func (s *Seeker) ByGUID(appGUID string) (string, error) {
	return appGUID, nil
}

//getAppGUID performs lookups against the CF API to convert org, space, and app
// names into the target app GUID
func (s *Seeker) getAppGUID(orgname, spacename, appname string) (guid string, err error) {
	log.Debugf("Getting org by name (%s) from CF API", orgname)
	org, err := s.CF.GetOrgByName(orgname)
	if err != nil {
		err = fmt.Errorf("While looking up given org: %s", err.Error())
		return
	}

	log.Debugf("Getting space by name (%s) and org GUID (%s) from CF API", spacename, org.Guid)
	space, err := s.CF.GetSpaceByName(spacename, org.Guid)
	if err != nil {
		err = fmt.Errorf("While looking up given space: %s", err.Error())
		return
	}

	log.Debugf("Getting app by name (%s), space GUID (%s), and org GUID (%s) from CF API", appname, space.Guid, org.Guid)
	app, err := s.CF.AppByName(appname, space.Guid, org.Guid)
	if err != nil {
		err = fmt.Errorf("While looking up given app: %s", err.Error())
		return
	}
	return app.Guid, nil
}

func canonizeIP(ip string) (canon string, err error) {
	intermediate := net.ParseIP(ip)
	if intermediate == nil {
		return "", fmt.Errorf("Could not interpret `%s` as IP address", ip)
	}
	return intermediate.String(), nil
}
