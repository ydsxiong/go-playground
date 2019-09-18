package database

import (
	"fmt"

	"github.com/Benchkram/errz"
	"github.com/jinzhu/gorm"
	"github.com/ydsxiong/go-playground/boxoffice/config"

	// need a sql driver accessing mysql db
	_ "github.com/go-sql-driver/mysql"
)

func NewGormDB(conf *config.Config) *gorm.DB {
	pwd := conf.DB.Password
	if pwd != "" {
		pwd = ":" + pwd
	}
	dbURI := fmt.Sprintf("%s%s@%s",
		conf.DB.Username,
		pwd,
		conf.DB.ConnectUri)

	gormdb, err := gorm.Open(conf.DB.Dialect, dbURI)
	errz.Fatal(err, "Could not connect database\n")
	defer errz.Recover(&err)

	return gormdb
}
