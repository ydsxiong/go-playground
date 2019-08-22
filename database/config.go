package database

type DBConfig struct {
	Dialect    string
	ConnectUri string
	Username   string
	Password   string
}

func GetConfig(dialect string, uri string, user string, password string) *DBConfig {
	return &DBConfig{
		Dialect:    dialect,
		ConnectUri: uri,
		Username:   user,
		Password:   password,
	}
}

/*
for mysql database on localhost
*/
func GetDefaultConfig() *DBConfig {
	return GetConfig(
		"mysql",
		"/localdb?charset=utf8&parseTime=True",
		"root",
		"")
}

func (c *DBConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type DBaux struct {
		Dialect    string `yaml:"dialect"`
		ConnectUri string `yaml:"connectUri"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
	}
	var aux struct {
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

func pushMapDataIntoconfig(c *DBConfig, data map[interface{}]interface{}) {
	db := data["database"].(map[interface{}]interface{})
	c.Dialect = getValueForStringKey(db, "dialect")
	c.ConnectUri = getValueForStringKey(db, "connectUri")
	c.Username = getValueForStringKey(db, "username")
	c.Password = getValueForStringKey(db, "password")
}

func getValueForStringKey(data map[interface{}]interface{}, key string) string {
	if v, ok := data[key].(string); ok {
		return v
	}
	return ""
}
