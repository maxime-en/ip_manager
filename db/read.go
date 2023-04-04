package db

import (
	"ip_manager/types"
)

// getHosts retrieves all hosts of a subnet from the Hosts table
func (db *MySQLDatabase) getHosts(subnet *types.Subnet) error {
	rows, err := db.Conn.Query(
		"SELECT HostKey, Address, Description FROM Hosts WHERE SubnetKey = ?",
		subnet.Key,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var strAddr string
		var description string
		err = rows.Scan(&key, &strAddr, &description)
		if err != nil {
			return err
		}
		err = subnet.NewHost(key, strAddr, description)
		if err != nil {
			return err
		}
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
		"SELECT SubnetKey, Cidr, Description FROM Subnets WHERE PrefixKey = ?",
		prefix.Key,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var strCidr string
		var description string
		err = rows.Scan(key, strCidr, description)
		if err != nil {
			return err
		}
		err = prefix.NewSubnet(key, strCidr, description)
		if err != nil {
			return err
		}
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

// getPrefixes retrieves all prefixes from the Prefixes table
func (db *MySQLDatabase) getPrefixes(svc *types.Service) error {
	rows, err := db.Conn.Query("SELECT PrefixKey, Cidr, Description FROM Prefixes")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var strCidr string
		var description string
		err = rows.Scan(&key, &strCidr, &description)
		if err != nil {
			return err
		}
		err = svc.NewPrefix(key, strCidr, description)
		if err != nil {
			return err
		}
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}
