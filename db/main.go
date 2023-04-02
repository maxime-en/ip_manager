package db

import (
	"database/sql"
	"fmt"
	"ip_manager/types"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLDatabase represents a connection to a MySQL database
type MySQLDatabase struct {
	Conn *sql.DB
}

// NewMySQLDatabase creates a new connection to a MySQL database
func NewMySQLDatabase(user, password, host, port, dbname string) (*MySQLDatabase, error) {
	// Creates the connection chain
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	// Opens it
	conn, err := sql.Open("mysql", uri)
	if err != nil {
		return nil, err
	}

	// Checks the conn
	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	// Returns it
	return &MySQLDatabase{
		Conn: conn,
	}, nil
}

// "Close" closes the connection
func (db *MySQLDatabase) Close() error {
	return db.Conn.Close()
}

// LoadService charge les IPs, Subnets et SubnetGroups en mémoire à partir de la base de données
func (db *MySQLDatabase) LoadService(service *types.Service) error {
	if err := db.createTablesIfNotExists(); err != nil {
		return err
	}

	if err := db.getPrefixes(service); err != nil {
		return err
	}

	// Charger les subnets et les ips pour tous les SubnetGroup
	for _, prefix := range service.Prefixes {
		err := db.getSubnets(prefix)
		if err != nil {
			return err
		}
		// Charger les ips pour tous les subnets du SubnetGroup actuel
		for _, subnet := range prefix.Subnets {
			err = db.getHosts(subnet)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Create database structure if it does not exist, called by LoadService
func (db *MySQLDatabase) createTablesIfNotExists() error {
	// Create Prefixes table
	if _, err := db.Conn.Exec(`
		CREATE TABLE IF NOT EXISTS Prefixes(
			PrefixKey VARCHAR(255) NOT NULL UNIQUE,
			Cidr VARCHAR(255) NOT NULL UNIQUE,
			Description TEXT(1000),
			Rfc1918compliant BOOLEAN NOT NULL DEFAULT true,
			PRIMARY KEY(PrefixKey)
		)
	`); err != nil {
		return err
	}

	// Create Subnets table
	if _, err := db.Conn.Exec(`
		CREATE TABLE IF NOT EXISTS Subnets (
			SubnetKey VARCHAR(255) NOT NULL UNIQUE, 
			Cidr VARCHAR(255) NOT NULL UNIQUE, 
			Description TEXT(1000), 
			PrefixKey VARCHAR(255) NOT NULL, 
			Rfc1918compliant BOOLEAN NOT NULL DEFAULT true, 
			PRIMARY KEY (SubnetKey), 
			FOREIGN KEY (PrefixKey) REFERENCES Prefixes(PrefixKey) ON DELETE CASCADE
		)
	`); err != nil {
		return err
	}

	// Create Hosts table
	if _, err := db.Conn.Exec(`
		CREATE TABLE IF NOT EXISTS Hosts (
			HostKey VARCHAR(255) NOT NULL UNIQUE, 
			Address VARCHAR(255) NOT NULL UNIQUE, 
			Description TEXT(1000), 
			SubnetKey VARCHAR(255) NOT NULL, 
			Rfc1918compliant BOOLEAN NOT NULL DEFAULT true, 
			PRIMARY KEY (HostKey), 
			FOREIGN KEY (SubnetKey) REFERENCES Subnets(SubnetKey) ON DELETE CASCADE
		)
	`); err != nil {
		return err
	}

	return nil
}
