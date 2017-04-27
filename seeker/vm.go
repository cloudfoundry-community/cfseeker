package seeker

import (
	"fmt"

	"github.com/cloudfoundry-community/gogobosh"
	"github.com/starkandwayne/goutils/log"
)

//VMInfo contains information about a VM...
type VMInfo struct {
	Name  string
	IP    string
	Index int
}

// GetVMWithIP searches the BOSH director for the VM with the IP you've given
func (s *Seeker) GetVMWithIP(ip string) (vm *VMInfo, err error) {
	log.Debugf("Getting VM with IP (%s)", ip)
	ip, err = canonizeIP(ip)
	if err != nil {
		return
	}
	s.acquireLock()
	defer s.releaseLock()
	tmpVM, found := s.vmdata[ip]
	if found {
		log.Debugf("Cache HIT for VM with IP (%s)", ip)
		vm = &tmpVM
		return
	}
	log.Debugf("Cache MISS for VM with IP (%s)", ip)
	//If we're here, we need to (try to) fetch the VM from BOSH
	err = s.cacheUntil(ip)
	if err != nil {
		err = fmt.Errorf("Error fetching VMs: %s", err.Error())
		return
	}

	log.Debugf("Checking for VM with IP (%s) after fetching data")
	tmpVM, found = s.vmdata[ip]
	if found {
		log.Debugf("Found VM with IP (%s) after fetching", ip)
		vm = &tmpVM
	} else {
		log.Debugf("Could not find VM with IP (%s)")
	}
	return
}

//This function will cycle through uncached deployments, storing their vm info
// until getting to the relevant deployment with the ip we need.
// It will only search in deployments that aren't already cached.
//
// SYNC: It is assumed that the lock is held by the caller of this function
func (s *Seeker) cacheUntil(ip string) (err error) {
	log.Debugf("Attempting fetch of VM with IP (%s)", ip)

	//If we've cached everything, then there's nothing to do here
	if len(s.cachedDeps) == len(s.config.BOSH.Deployments) {
		log.Debugf("Aborting fetch for ip (%s): Nothing left to fetch", ip)
		return
	}

	ip, err = canonizeIP(ip)
	if err != nil {
		return
	}

	// For this to work, deployments cached must be done in the same order as the
	// full deployment list. Basically, the loop continues where the cache left
	// off
	for _, depName := range s.config.BOSH.Deployments[len(s.cachedDeps):] {
		var vms []gogobosh.VM
		//Go get the VMs in this particular deployment
		log.Debugf("Contacting BOSH Director for VMs in deployment with name (%s)", depName)
		vms, err = s.bosh.GetDeploymentVMs(depName)
		if err != nil {
			return fmt.Errorf("Error while getting VMs for deployment `%s`: %s", depName, err.Error())
		}

		log.Debugf("Inserting VMs into local memory cache")
		//Populate the cache with the VMs we got
		for _, vm := range vms {
			vm.IPs[0], err = canonizeIP(vm.IPs[0])
			if err != nil {
				return
			}
			s.vmdata[vm.IPs[0]] = VMInfo{
				Name:  vm.JobName,
				IP:    vm.IPs[0],
				Index: vm.Index,
			}
		}
		log.Debugf("Cached %d VMs", len(vms))
		//Bail out if we got our target ip
		if _, found := s.vmdata[ip]; found {
			log.Debugf("Found target IP (%s) in deployment (%s)", ip, depName)
			return
		}
	}
	log.Debugf("Fetched all deployments but didn't find IP (%s)", ip)
	return
}

func (s *Seeker) invalidateVMs() {
	log.Debugf("Invalidating cache for Seeker (%p)", s)
	s.acquireLock()
	defer s.releaseLock()
	s.vmdata = map[string]VMInfo{}
	s.cachedDeps = []string{}
	log.Debugf("Cache invalidated for Seeker (%p)", s)
}

func (s *Seeker) acquireLock() {
	log.Debugf("Acquiring lock for Seeker (%p)", s)
	s.lock.Lock()
}

func (s *Seeker) releaseLock() {
	log.Debugf("Releasing lock for Seeker (%p)", s)
	s.lock.Unlock()
}
