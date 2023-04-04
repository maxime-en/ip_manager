package types

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

type Prefix struct {
	Key          string             `json:"key"`
	Cidr         *net.IPNet         `json:"-"`
	StrCidr      string             `json:"cidr"`
	Description  string             `json:"description"`
	Service      *Service           `json:"-"` // link to Service
	Subnets      []*Subnet          `json:"-"` // link to Subnets
	SubnetsByKey map[string]*Subnet `json:"-"`
}

const (
	minGroupMask = 23
)

// Prefix generator
func NewPrefix(key string, strCidr string, service *Service, description string) (*Prefix, error) {
	// Check if service exists
	if service == nil {
		return nil, fmt.Errorf("The service doesn't exist, unable to create a prefix.")
	}

	// Check if the key contains only allowed chars
	if !isValidKey(key) {
		return nil, fmt.Errorf("Invalid key: %s. Allowed chars: [a-z0-9.-]", key)
	}

	// Check if the key is not already used
	if _, ok := service.PrefixesByKey[key]; ok {
		return nil, fmt.Errorf("Already used key: %s.", key)
	}

	// Parse CIDR from string
	ip, oCidr, err := net.ParseCIDR(strCidr)
	if err != nil {
		return nil, fmt.Errorf("Invalid prefix CIDR: %s", strCidr)
	}

	// Mask must be bigger than minGroupMask
	oCidrMaskSize, _ := oCidr.Mask.Size()
	if oCidrMaskSize > minGroupMask {
		return nil, fmt.Errorf("Invalid prefix mask: %s. It must be bigger than %s.", strCidr, strconv.Itoa(minGroupMask))
	}

	// Is CIDR unique ?
	for _, existingPrefix := range service.Prefixes {
		if existingPrefix.Cidr.Contains(ip) {
			return nil, fmt.Errorf("The CIDR %s is already used by prefix %s.", strCidr, existingPrefix.Key)
		}
	}

	// Create fields
	prefix := &Prefix{
		Key:          key,
		Cidr:         oCidr,
		StrCidr:      strCidr,
		Service:      service,
		Description:  description,
		Subnets:      make([]*Subnet, 0),
		SubnetsByKey: make(map[string]*Subnet),
	}

	// Link the prefix to the service
	service.Prefixes = append(service.Prefixes, prefix)
	service.PrefixesByKey[key] = prefix

	return prefix, nil
}

// Method to call subnet generator
func (prefix *Prefix) NewSubnet(key string, strCidr string, description string) error {
	if _, err := NewSubnet(key, strCidr, prefix, description); err != nil {
		return err
	}
	return nil
}

// Print in json format a prefix
func (prefix *Prefix) ToJSON() ([]byte, error) {
	if prefix == nil {
		return nil, fmt.Errorf("The prefix doesn't exist, unable to print it.")
	} else {
		return json.Marshal(prefix)
	}
}

// Evaluate if the prefix is RFC1819 compliant
func (prefix *Prefix) IsPrivate() (bool, error) {
	if prefix != nil {
		return prefix.Cidr.IP.IsPrivate(), nil
	}
	return false, fmt.Errorf("The prefix does not exists, unable to evaluate it.")
}

// Edit an existing Prefix
func (prefix *Prefix) Modify(strCidr string, description string) error {
	// Existence check
	if prefix == nil {
		return fmt.Errorf("The prefix doesn't exist, unable to modify it.")
	}
	// Empty check
	if len(prefix.Subnets) != 0 {
		return fmt.Errorf("Unable to modify the prefix %s: it contains subnet(s).", prefix.Key)
	}

	// Parse CIDR from string
	ip, oCidr, err := net.ParseCIDR(strCidr)
	if err != nil {
		return fmt.Errorf("Invalid prefix CIDR : %s", strCidr)
	}

	// Mask must be bigger than minGroupMask
	oCidrMaskSize, _ := oCidr.Mask.Size()
	if oCidrMaskSize > minGroupMask {
		return fmt.Errorf("Invalid mask: %s. It must be bigger than %s.", strCidr, strconv.Itoa(minGroupMask))
	}

	// Is CIDR unique ?
	for _, existingPrefix := range prefix.Service.Prefixes {
		if existingPrefix.Cidr.Contains(ip) && existingPrefix.Key != existingPrefix.Key {
			return fmt.Errorf("The CIDR %s is already used by %s.", strCidr, existingPrefix.Key)
		}
	}

	// Update fields
	prefix.Cidr = oCidr
	prefix.StrCidr = strCidr
	prefix.Description = description

	return nil
}

// Delete an existing prefix
func (prefix *Prefix) Delete() error {
	// Existence check
	if prefix == nil {
		return fmt.Errorf("The prefix doesn't exist, unable to delete it.")
	}

	// Delete linked subnets
	for _, subnet := range prefix.Subnets {
		err := subnet.Delete()
		if err != nil {
			return err
		}
	}

	// Delete prefix from service
	for i, existingPrefix := range prefix.Service.Prefixes {
		if existingPrefix == prefix {
			prefix.Service.Prefixes = append(prefix.Service.Prefixes[:i], prefix.Service.Prefixes[i+1:]...)
			delete(prefix.Service.PrefixesByKey, prefix.Key)
			return nil
		}
	}

	return fmt.Errorf("Unable to delete the prefix: can't find it.")
}
