package types

import (
	"fmt"
	"regexp"
)

// Contains all prefixes of a service
type Service struct {
	Prefixes      []*Prefix
	PrefixesByKey map[string]*Prefix
	Subnets       []*Subnet
	SubnetsByKey  map[string]*Subnet
	Hosts         []*Host
	HostsByKey    map[string]*Host
}

// Prefix generator
func NewService() *Service {
	return &Service{
		Prefixes:      make([]*Prefix, 0),
		PrefixesByKey: make(map[string]*Prefix),
		Subnets:       make([]*Subnet, 0),
		SubnetsByKey:  make(map[string]*Subnet),
		Hosts:         make([]*Host, 0),
		HostsByKey:    make(map[string]*Host),
	}
}

// Method to call prefix generator
func (svc *Service) NewPrefix(key string, strCidr string, description string) error {
	if _, err := NewPrefix(key, strCidr, svc, description); err != nil {
		return err
	}
	return nil
}

// Delete an existing service
func (svc *Service) Delete() error {
	// Existence check
	if svc == nil {
		return fmt.Errorf("The service doesn't exist, unable to delete it.")
	}

	// Delete linked Subnets
	for _, prefix := range svc.Prefixes {
		err := prefix.Delete()
		if err != nil {
			return err
		}
	}

	// Clear the map and slice
	svc.PrefixesByKey = make(map[string]*Prefix)
	svc.SubnetsByKey = make(map[string]*Subnet)
	svc.HostsByKey = make(map[string]*Host)
	svc.Prefixes = make([]*Prefix, 0)
	svc.Subnets = make([]*Subnet, 0)
	svc.Hosts = make([]*Host, 0)

	return nil
}

// Check if key only contains allowed chars and is not longer than 50 chars
func isValidKey(key string) bool {
	// 50 chars max
	if len(key) > 50 {
		return false
	}

	// Allowed chars: [a-z0-9.-]
	match, _ := regexp.MatchString("^[a-z0-9.-]+$", key)
	return match
}
