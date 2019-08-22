package database

import (
	"fmt"

	"github.com/Benchkram/errz"
	"github.com/jinzhu/gorm"
)

func InitDB(config *DBConfig) *gorm.DB {
	pwd := config.Password
	if pwd != "" {
		pwd = ":" + pwd
	}
	dbURI := fmt.Sprintf("%s%s@%s",
		config.Username,
		pwd,
		config.ConnectUri)

	db, err := gorm.Open(config.Dialect, dbURI)
	errz.Fatal(err, "Could not connect database\n")
	defer errz.Recover(&err)
	// if err != nil {
	// 	log.Fatal("Could not connect database")
	// }
	return db
}

func InitDBLocal() *gorm.DB {
	return InitDB(GetDefaultConfig())
}
