package seeker

import (
	"fmt"

	"github.com/cloudfoundry-community/gogobosh"
)

//VMInfo contains information about a VM...
type VMInfo struct {
	Name  string
	IP    string
	Index int
}

// GetVMWithIP searches the BOSH director for the VM with the IP you've given
func (s *Seeker) GetVMWithIP(ip string) (vm *VMInfo, err error) {
	ip, err = canonizeIP(ip)
	if err != nil {
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	tmpVM, found := s.vmdata[ip]
	if found {
		vm = &tmpVM
		return
	}
	//If we're here, we need to (try to) fetch the VM from BOSH
	err = s.cacheUntil(ip)
	if err != nil {
		err = fmt.Errorf("Error fetching VMs: %s", err.Error())
		return
	}

	tmpVM, found = s.vmdata[ip]
	if found {
		vm = &tmpVM
	}
	return
}

//This function will cycle through uncached deployments, storing their vm info
// until getting to the relevant deployment with the ip we need.
// It will only search in deployments that aren't already cached.
//
// SYNC: It is assumed that the lock is held by the caller of this function
func (s *Seeker) cacheUntil(ip string) (err error) {
	//If we've cached everything, then there's nothing to do here
	if len(s.cachedDeps) == len(s.config.BOSH.Deployments) {
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
		vms, err = s.bosh.GetDeploymentVMs(depName)
		if err != nil {
			return fmt.Errorf("Error while getting VMs for deployment `%s`: %s", depName, err.Error())
		}
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
		//Bail out if we got our target ip
		if _, found := s.vmdata[ip]; found {
			return
		}
	}
	return
}

func (s *Seeker) invalidateVMs() {
	s.lock.Lock()
	s.vmdata = map[string]VMInfo{}
	s.cachedDeps = []string{}
	s.lock.Unlock()
}
