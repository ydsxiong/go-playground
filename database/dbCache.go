package database

import (
	"database/sql"

	// need a sql driver for the db fetching
	_ "github.com/go-sql-driver/mysql"
)

type dbWrap struct {
	db  *sql.DB
	err *error
}

// set up a channel for requesting and receiving a sql db from the cache
var dbFetchChannel chan dbWrap

func initDBPool() {
	if dbFetchChannel != nil {
		return
	}

	dbFetchChannel = make(chan dbWrap)
	// fire up a go routine to create and cache a sql db there, using the fetch channel above
	go func(dbChan chan dbWrap) {
		var cachedDB *dbWrap
		for {
			select {
			case <-dbChan: // check for a request for a sql db
				if cachedDB == nil || *(cachedDB.err) != nil || cachedDB.db.Ping() != nil {
					db, err := sql.Open("mysql", "root@/localdb")
					cachedDB = &dbWrap{db, &err}
				}
				dbChan <- *cachedDB // respond with the cached sql db
			}
		}
	}(dbFetchChannel)
}

// FetchDBFromCache for get a sql db object in cache
func FetchDBFromCache() (db *sql.DB, err *error) {
	initDBPool()
	dbFetchChannel <- dbWrap{nil, nil} // request a sql db off the cache
	cachedDB := <-dbFetchChannel       // receive the sql db sent from the cache
	return cachedDB.db, cachedDB.err
}

// CloseDB to close the sql object if needed, however unnecessarily it may be
func CloseDB() {
	db, _ := FetchDBFromCache()
	db.Close()
}
