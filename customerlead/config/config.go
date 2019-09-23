package config

type Config struct {
	ServerPort string
	DB         *DBConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Dialect    string
	ConnectUri string
	Username   string
	Password   string
}

func GetConfig(dialect string, uri string, user string, password string) *Config {
	return &Config{
		DB: &DBConfig{
			Dialect:    dialect,
			ConnectUri: uri,
			Username:   user,
			Password:   password,
		},
	}
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type Serveraux struct {
		Port string `yaml:"port"`
	}
	type DBaux struct {
		Dialect    string `yaml:"dialect"`
		ConnectUri string `yaml:"connectUri"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
	}
	var aux struct {
		Serveraux `yaml:"server"`
		DBaux     `yaml:"database"`
	}

	err := unmarshal(&aux)
	if err != nil {
		return err
	}

	c.ServerPort = aux.Port
	c.DB.Dialect = aux.Dialect
	c.DB.ConnectUri = aux.ConnectUri
	c.DB.Username = aux.Username
	c.DB.Password = aux.Password
	return nil
}
