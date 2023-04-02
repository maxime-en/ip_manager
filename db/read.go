package db

import (
	"ip_manager/types"
	"net"
)

// getHosts retrieves all hosts of a subnet from the Hosts table
func (db *MySQLDatabase) getHosts(subnet *types.Subnet) error {
	rows, err := db.Conn.Query(
		"SELECT HostKey, Address, Description, SubnetKey, Rfc1918compliant FROM Hosts WHERE SubnetKey = ?",
		subnet.Key,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var host types.Host
		var StrAddr string
		err = rows.Scan(&host.Key, &StrAddr, &host.Description, &host.Subnet.Key, &host.Rfc1918compliant)
		if err != nil {
			return err
		}
		host.Address = net.ParseIP(StrAddr)
		subnet.Hosts = append(subnet.Hosts, &host)
		subnet.HostsByKey[host.Key] = &host
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

// getSubnets retrieves all subnets of a prefix from the Subnets table
func (db *MySQLDatabase) getSubnets(prefix *types.Prefix) error {

	rows, err := db.Conn.Query(
		"SELECT SubnetKey, Cidr, Description, PrefixKey, Rfc1918compliant FROM Subnets WHERE PrefixKey = ?",
		prefix.Key,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var subnet types.Subnet
		var StrCidr string
		err = rows.Scan(&subnet.Key, &StrCidr, &subnet.Description, &subnet.Prefix.Key, &subnet.Rfc1918compliant)
		if err != nil {
			return err
		}
		_, subnet.Cidr, err = net.ParseCIDR(StrCidr)
		if err != nil {
			return err
		}
		prefix.Subnets = append(prefix.Subnets, &subnet)
		prefix.SubnetsByKey[subnet.Key] = &subnet
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

// getPrefixes retrieves all prefixes from the Prefixes table
func (db *MySQLDatabase) getPrefixes(svc *types.Service) error {
	rows, err := db.Conn.Query("SELECT PrefixKey, Cidr, Description, Rfc1918compliant FROM Prefixes")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var prefix types.Prefix
		var StrCidr string
		err = rows.Scan(&prefix.Key, &StrCidr, &prefix.Description, &prefix.Rfc1918compliant)
		if err != nil {
			return err
		}
		_, prefix.Cidr, err = net.ParseCIDR(StrCidr)
		if err != nil {
			return err
		}
		svc.Prefixes = append(svc.Prefixes, &prefix)
		svc.PrefixesByKey[prefix.Key] = &prefix
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}
