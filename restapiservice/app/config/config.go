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

/*
for mysql database on localhost
*/
func GetDefaultConfig() *Config {
	return GetConfig(
		"mysql",
		"/localdb?charset=utf8&parseTime=True",
		"root",
		"")
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type Serveraux struct {
		Port string `yaml:"port"`
	}
	type DBaux struct {
		Dialect string `yaml:"dialect"`
		ConnectUri string `yaml:"connectUri"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
	var aux struct {
		Serveraux `yaml:"server"`
		 DBaux `yaml:"database"`
	}

	// map can work here smoothly
	//data := make(map[interface{}]interface{})
	err := unmarshal(&aux)
	if err != nil {
		return err
	}

	//pushMapDataIntoconfig(c, data)

	return nil

}

func pushMapDataIntoconfig(c *Config, data map[interface{}]interface{}) {
	server := data["server"].(map[interface{}]interface{})
	c.ServerPort = getValueForStringKey(server, "port")

	db := data["database"].(map[interface{}]interface{})
	c.DB.Dialect = getValueForStringKey(db, "dialect")
	c.DB.ConnectUri = getValueForStringKey(db, "connectUri")
	c.DB.Username = getValueForStringKey(db, "username")
	c.DB.Password = getValueForStringKey(db, "password")
}

func getValueForStringKey(data map[interface{}]interface{}, key string) string {
	if v, ok := data[key].(string); ok {
		return v
	}
	return ""
}
