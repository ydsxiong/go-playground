package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ydsxiong/playground/customerlead/config"
	"github.com/ydsxiong/playground/customerlead/database"
	"github.com/ydsxiong/playground/customerlead/lead"
	"gopkg.in/yaml.v2"
)

const (
	ENV_DB_DIALECT     string = "DB_DIALECT"
	ENV_DB_CONNECT_URI string = "DB_CONNECT_URI"
	ENV_DB_USERNAME    string = "DB_USERNAME"
	EVN_DB_PASSWORD    string = "DB_PASSWORD"
	ENV_PORT           string = "PORT"

	USE_IN_MEMORY_STORE string = "in-memory-store"
	USE_DATABASE_STORE  string = "database-store"
	USE_FILE_STORE      string = "file-system-store"

	DATA_STORE_FILE = "../../customer-leads.db.json"
)

func main() {

	//
	//these various data storage setups are just for demo purpose, not really needed for running the app here
	//
	dataStoreType := flag.String("datasource", USE_FILE_STORE, "data store source")
	flag.Parse()

	var closeStore *func()

	var datastore lead.LeadStore
	if *dataStoreType == USE_IN_MEMORY_STORE {
		datastore = setupInMemoryStore()
	} else if *dataStoreType == USE_DATABASE_STORE {
		datastore = setupDatabaseStore()
	} else {
		datastore, closeStore = setupFileSystemStore()
	}

	if closeStore != nil {
		defer (*closeStore)()
	}

	/////////////////////////////////////////////////////////////////////
	// the main code is here: set up the webserver and start it up
	//
	server, err := lead.NewLeadServer(datastore)
	if err != nil {
		log.Fatalf("Problem with setting up the server, %v", err)
	}

	// port number can be configured from ENV or bootstap config file
	if http.ListenAndServe(":9000", server) != nil {
		log.Fatalf("Couldn't listen to 9000 port, %v", err)
	}
}

/**
these various data storage setups are just for demo purpose, not really needed for running the app here
*/
func setupInMemoryStore() lead.LeadStore {
	return lead.NewInMemoryDataStore()
}

func setupFileSystemStore() (lead.LeadStore, *func()) {
	store, closeStore, err := lead.LoadUpFileStore(DATA_STORE_FILE)
	if err != nil {
		log.Fatalf("Problem with loading in file store, %v", err)
	}
	return store, &closeStore
}

func setupDatabaseStore() lead.LeadStore {
	conf := config.GetConfig(
		os.Getenv(ENV_DB_DIALECT),
		os.Getenv(ENV_DB_CONNECT_URI),
		os.Getenv(ENV_DB_USERNAME),
		os.Getenv(EVN_DB_PASSWORD))
	conf.ServerPort = os.Getenv(ENV_PORT)

	configPtr := flag.String("config", "config.local.yml", "a config file path")

	configdata, err := ioutil.ReadFile(*configPtr)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(configdata, conf); err != nil {
		log.Fatal(err)
	}

	gormDb := database.NewGormDB(conf)
	// auto create customer leads table if not existed
	gormDb.AutoMigrate(&lead.Lead{})

	return lead.NewDatabaseStore(gormDb)
}
