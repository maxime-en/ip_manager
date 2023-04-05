package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// Config "main"
	DbSaveInterval int  `yaml:"db_save_interval"`
	Debug          bool `yaml:"debug"`
	// Config BDD
	MySQLHost     string `yaml:"mysql_host"`
	MySQLPort     string `yaml:"mysql_port"`
	MySQLUser     string `yaml:"mysql_user"`
	MySQLPassword string `yaml:"mysql_password"`
	MySQLDatabase string `yaml:"mysql_database"`
	// API config
	Listen           string   `yaml:"listen"`
	AuthorizedTokens []string `yaml:"authorized_tokens"`
}

func LoadConfig() (*Config, error) {
	configFile := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	file, err := os.Open(*configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
