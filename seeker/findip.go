package seeker

import "fmt"

//FindOutput contains the return values from a call to Find()
type FindOutput struct {
	VMName string
	Host   string
}

// FindIPs takes the App GUID given and queries the CF API to get the IP
// addresses of the VMs on which the instances of the application are located
// and returns those. If inputErr is non-nil, the function will bail early with
// the given error.
func (s *Seeker) FindIPs(guid string, inputErr error) (ips []string, err error) {
	if inputErr != nil {
		return nil, inputErr
	}

	return s.getAppHosts(guid)
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
	org, err := s.client.GetOrgByName(orgname)
	if err != nil {
		err = fmt.Errorf("While looking up given org: %s", err.Error())
		return
	}

	space, err := s.client.GetSpaceByName(spacename, org.Guid)
	if err != nil {
		err = fmt.Errorf("While looking up given space: %s", err.Error())
		return
	}

	app, err := s.client.AppByName(appname, space.Guid, org.Guid)
	if err != nil {
		err = fmt.Errorf("While looking up given app: %s", err.Error())
		return
	}
	return app.Guid, nil
}

func (s *Seeker) getAppHosts(guid string) (ips []string, err error) {
	statsMap, err := s.client.GetAppStats(guid)
	if err != nil {
		err = fmt.Errorf("Error when getting stats for app with GUID `%s` (Is the app running?)", guid)
		return
	}
	if len(statsMap) == 0 {
		err = fmt.Errorf("No stats found for app with GUID `%s`", guid)
		return
	}

	for _, stats := range statsMap {
		ips = append(ips, stats.Stats.Host)
	}
	return
}
