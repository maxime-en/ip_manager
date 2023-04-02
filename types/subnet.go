package types

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

type Subnet struct {
	Key              string           `json:"key"`
	Cidr             *net.IPNet       `json:"-"`
	StrCidr          string           `json:"cidr"`
	Description      string           `json:"description"`
	Rfc1918compliant bool             `json:"rfc1918compliant"`
	Prefix           *Prefix          `json:"SubnetGroup"` // link to SubnetGroup
	Hosts            []*Host          `json:"-"`           // link to IPs
	HostsByKey       map[string]*Host `json:"-"`
}

const (
	minSubMask = 24
)

// Subnet generator
func NewSubnet(key string, strCidr string, prefix *Prefix, description string) (*Subnet, error) {
	// Check if prefix exists
	if prefix == nil {
		return nil, fmt.Errorf("The prefix doesn't exist, unable to create a subnet.")
	}
	// Check if the key contains only allowed chars
	if !isValidKey(key) {
		return nil, fmt.Errorf("Invalid key: %s. Allowed chars: [a-z0-9.-]", key)
	}

	// Check if the key is not already used
	if _, ok := prefix.SubnetsByKey[key]; ok {
		return nil, fmt.Errorf("Already used key: %s.", key)
	}

	// Parse CIDR from string
	ip, oCidr, err := net.ParseCIDR(strCidr)
	if err != nil {
		return nil, fmt.Errorf("Invalid subnet CIDR: %s", strCidr)
	}

	// Mask must be bigger than minSubMask
	oCidrMaskSize, _ := oCidr.Mask.Size()
	if oCidrMaskSize > minSubMask {
		return nil, fmt.Errorf("Invalid subnet mask: %s. It must be bigger than %s.", strCidr, strconv.Itoa(minSubMask))
	}

	// CIDR have to be contained in the Prefix
	if !prefix.Cidr.Contains(ip) {
		return nil, fmt.Errorf("Invalid subnet CIDR: %s. It must be contained in the prefix %s.", strCidr, prefix.Cidr)
	}

	// Is CIDR unique ?
	for _, existingSubnet := range prefix.Subnets {
		if existingSubnet.Cidr.Contains(ip) {
			return nil, fmt.Errorf("The CIDR %s is already used by subnet %s.", strCidr, existingSubnet.Key)
		}
	}

	// Create fields
	subnet := &Subnet{
		Key:              key,
		Cidr:             oCidr,
		StrCidr:          strCidr,
		Prefix:           prefix,
		Description:      description,
		Rfc1918compliant: ip.IsPrivate(),
		Hosts:            make([]*Host, 0),
		HostsByKey:       make(map[string]*Host),
	}

	// Link the subnet to the parent prefix
	prefix.Subnets = append(prefix.Subnets, subnet)
	prefix.SubnetsByKey[key] = subnet

	return subnet, nil
}

// Print in json format a subnet
func (subnet *Subnet) ToJSON() ([]byte, error) {
	if subnet == nil {
		return nil, fmt.Errorf("The subnet doesn't exist, unable to print it.")
	} else {
		return json.Marshal(subnet)
	}
}

// Edit an existing Subnet
func (subnet *Subnet) Modify(strCidr string, description string) error {
	// Existence check
	if subnet == nil {
		return fmt.Errorf("The subnet doesn't exist, unable to modify it.")
	}
	// Empty check
	if len(subnet.Hosts) != 0 {
		return fmt.Errorf("Unable to modify the subnet %s: it contains host(s).", subnet.Key)
	}

	// Parse CIDR from string
	ip, oCidr, err := net.ParseCIDR(strCidr)
	if err != nil {
		return fmt.Errorf("Invalid CIDR: %s", strCidr)
	}

	// Mask must be bigger than minSubMask
	oCidrMaskSize, _ := oCidr.Mask.Size()
	if oCidrMaskSize > minSubMask {
		return fmt.Errorf("Invalid mask: %s. It must be bigger than %s.", strCidr, strconv.Itoa(minSubMask))
	}

	// CIDR have to be contained in the SubnetGroup
	if !subnet.Prefix.Cidr.Contains(ip) {
		return fmt.Errorf("Invalid CIDR: %s. It must be contained in the prefix %s.", strCidr, subnet.Prefix.Cidr)
	}

	// Is CIDR unique ?
	for _, existingSubnet := range subnet.Prefix.Subnets {
		if existingSubnet.Key != subnet.Key && existingSubnet.Cidr.Contains(ip) {
			return fmt.Errorf("The CIDR %s is already used by subnet %s.", strCidr, existingSubnet.Key)
		}
	}

	// Update fields
	subnet.Cidr = oCidr
	subnet.StrCidr = strCidr
	subnet.Description = description
	subnet.Rfc1918compliant = ip.IsPrivate()

	return nil
}

// Delete an existing subnet
func (subnet *Subnet) Delete() error {
	// Existence check
	if subnet == nil {
		return fmt.Errorf("The subnet doesn't exist, unable to delete it.")
	}

	// Delete linked hosts
	for _, host := range subnet.Hosts {
		err := host.Delete()
		if err != nil {
			return err
		}
	}

	// Delete subnet from parent prefix
	for i, existingSubnet := range subnet.Prefix.Subnets {
		if existingSubnet == subnet {
			subnet.Prefix.Subnets = append(subnet.Prefix.Subnets[:i], subnet.Prefix.Subnets[i+1:]...)
			delete(subnet.Prefix.SubnetsByKey, subnet.Key)
			return nil
		}
	}

	return fmt.Errorf("Unable to delete the subnet: can't find it.")
}
