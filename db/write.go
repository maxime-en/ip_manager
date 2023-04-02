package db

import "ip_manager/types"

// AddHost adds an host record to the SQL table
func (db *MySQLDatabase) AddHost(host *types.Host) error {
	_, err := db.Conn.Exec(
		"INSERT INTO Hosts (HostKey, Address, Description, SubnetKey, Rfc1918compliant) VALUES (?, ?, ?, ?, ?)",
		host.Key,
		host.StrAddress,
		host.Description,
		host.Subnet.Key,
		host.Rfc1918compliant,
	)
	if err != nil {
		return err
	}
	return nil
}

// UpdateHost updates a host record to the SQL table
func (db *MySQLDatabase) UpdateHost(host *types.Host) error {
	_, err := db.Conn.Exec(
		"UPDATE Hosts SET Address = ? , Description = ? , Rfc1918compliant = ? WHERE HostKey = ?",
		host.StrAddress,
		host.Description,
		host.Rfc1918compliant,
		host.Key,
	)
	if err != nil {
		return err
	}
	return nil
}

// DeleteHost deletes a SQL record to the SQL table
func (db *MySQLDatabase) DeleteHost(host *types.Host) error {
	_, err := db.Conn.Exec(
		"DELETE FROM Hosts WHERE HostKey = ?",
		host.Key,
	)
	if err != nil {
		return err
	}
	return nil
}

// AddSubnet adds an subnet record to the SQL table
func (db *MySQLDatabase) AddSubnet(subnet *types.Subnet) error {
	_, err := db.Conn.Exec(
		"INSERT INTO Subnets (SubnetKey, Cidr, Description, PrefixKey, Rfc1918Compliant) VALUES (?, ?, ?, ?, ?)",
		subnet.Key,
		subnet.StrCidr,
		subnet.Description,
		subnet.Prefix.Key,
		subnet.Rfc1918compliant,
	)
	if err != nil {
		return err
	}
	return nil
}

// UpdateSubnet updates a host record to the SQL table
func (db *MySQLDatabase) UpdateSubnet(subnet *types.Subnet) error {
	_, err := db.Conn.Exec(
		"UPDATE Subnets SET Cidr = ? , Description = ? , Rfc1918compliant = ? WHERE SubnetKey = ?",
		subnet.StrCidr,
		subnet.Description,
		subnet.Rfc1918compliant,
		subnet.Key,
	)
	if err != nil {
		return err
	}
	return nil
}

// DeleteSubnet deletes a SQL record to the SQL table
func (db *MySQLDatabase) DeleteSubnet(subnet *types.Subnet) error {
	_, err := db.Conn.Exec(
		"DELETE FROM Subnets WHERE SubnetKey = ?",
		subnet.Key,
	)
	if err != nil {
		return err
	}
	return nil
}

// AddPrefix adds an subnet record to the SQL table
func (db *MySQLDatabase) AddPrefix(prefix *types.Prefix) error {
	_, err := db.Conn.Exec(
		"INSERT INTO Prefixes (PrefixKey, Cidr, Description, Rfc1918compliant) VALUES (?, ?, ?, ?)",
		prefix.Key,
		prefix.StrCidr,
		prefix.Description,
		prefix.Rfc1918compliant,
	)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePrefix updates a host record to the SQL table
func (db *MySQLDatabase) UpdatePrefix(prefix *types.Prefix) error {
	_, err := db.Conn.Exec(
		"UPDATE Prefixes SET Cidr = ? , Description = ? , Rfc1918compliant = ? WHERE PrefixKey = ?",
		prefix.StrCidr,
		prefix.Description,
		prefix.Rfc1918compliant,
		prefix.Key,
	)
	if err != nil {
		return err
	}
	return nil
}

// DeletePrefix deletes a SQL record to the SQL table
func (db *MySQLDatabase) DeletePrefix(prefix *types.Prefix) error {
	_, err := db.Conn.Exec(
		"DELETE FROM Prefixes WHERE PrefixKey = ?",
		prefix.Key,
	)
	if err != nil {
		return err
	}
	return nil
}
