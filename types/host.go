package types

import (
	"encoding/json"
	"fmt"
	"net"
)

type Host struct {
	Key              string  `json:"key"`
	Address          net.IP  `json:"-"`
	StrAddress       string  `json:"address"`
	Description      string  `json:"description"`
	Rfc1918compliant bool    `json:"rfc1918compliant"`
	Subnet           *Subnet `json:"Subnet"` // link to subnet
}

// Host generator
func NewHost(key string, strAddr string, subnet *Subnet, description string) (*Host, error) {
	// Check if subnet exists
	if subnet == nil {
		return nil, fmt.Errorf("The subnet doesn't exist, unable to create host.")
	}

	// Check if the key contains only allowed chars
	if !isValidKey(key) {
		return nil, fmt.Errorf("Invalid key: %s. Allowed chars: [a-z0-9.-]", key)
	}

	// Check if the key is not already used
	if _, ok := subnet.HostsByKey[key]; ok {
		return nil, fmt.Errorf("Already used key: %s", key)
	}

	// Parse host IP from string
	oAddr := net.ParseIP(strAddr)
	if oAddr == nil {
		return nil, fmt.Errorf("Invalid IP: %s", strAddr)
	}

	// CIDR have to be contained in the SubnetGroup
	if !subnet.Cidr.Contains(oAddr) {
		return nil, fmt.Errorf("Invalid IP: %s. It must be contained in the subnet (%s).", oAddr, subnet.Cidr)
	}

	// Is IP addr unique ?
	for _, existingIP := range subnet.Hosts {
		if existingIP.Address.Equal(oAddr) {
			return nil, fmt.Errorf("The IP %s is already used by host %s.", strAddr, existingIP.Key)
		}
	}

	// Create fields
	host := &Host{
		Key:              key,
		Address:          oAddr,
		StrAddress:       strAddr,
		Subnet:           subnet,
		Description:      description,
		Rfc1918compliant: oAddr.IsPrivate(),
	}

	// Link the IP to the Subnet
	subnet.Hosts = append(subnet.Hosts, host)
	subnet.HostsByKey[key] = host

	return host, nil
}

// Print in json format a SubnetGroup
func (ip *Host) ToJSON() ([]byte, error) {
	if ip == nil {
		return nil, fmt.Errorf("The IP doesn't exist, unable to print it.")
	} else {
		return json.Marshal(ip)
	}
}

// Edit an existing host
func (ip *Host) Modify(strAddr string, description string) error {
	if ip == nil {
		return fmt.Errorf("The Host doesn't exist, unable to edit it.")
	}
	// Parse host IP from string
	oAddr := net.ParseIP(strAddr)
	if oAddr == nil {
		return fmt.Errorf("Invalid IP: %s", strAddr)
	}

	// Is IP addr unique ?
	for _, existingIP := range ip.Subnet.Hosts {
		if existingIP.Address.Equal(oAddr) && existingIP.Key != ip.Key {
			return fmt.Errorf("L'adresse IP %s est déjà utilisée par %s", strAddr, existingIP.Key)
		}
	}

	// Update fields
	ip.Address = oAddr
	ip.StrAddress = strAddr
	ip.Description = description
	ip.Rfc1918compliant = oAddr.IsPrivate()

	return nil
}

// Delete an existing host
func (host *Host) Delete() error {
	if host == nil {
		return fmt.Errorf("The host doesn't exist, unable to delete it.")
	}
	// Delete host from parent subnet
	for i, existingHost := range host.Subnet.Hosts {
		if existingHost == host {
			host.Subnet.Hosts = append(host.Subnet.Hosts[:i], host.Subnet.Hosts[i+1:]...)
			delete(host.Subnet.HostsByKey, host.Key)
			return nil
		}
	}

	return fmt.Errorf("Unable to delete the host: can't find it.")
}
