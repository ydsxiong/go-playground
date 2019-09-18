package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	glog "github.com/go-kit/kit/log"
	_ "github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ydsxiong/go-playground/boxoffice/config"
	"github.com/ydsxiong/go-playground/boxoffice/controller"
	"github.com/ydsxiong/go-playground/boxoffice/database"
	"github.com/ydsxiong/go-playground/boxoffice/model"
	"github.com/ydsxiong/go-playground/boxoffice/services/inprogress"
	"github.com/ydsxiong/go-playground/boxoffice/services/registeredguest"
	"gopkg.in/yaml.v2"
)

const (
	ENV_DB_DIALECT          string = "DB_DIALECT"
	ENV_DB_CONNECT_URI      string = "DB_CONNECT_URI"
	ENV_DB_USERNAME         string = "DB_USERNAME"
	EVN_DB_PASSWORD         string = "DB_PASSWORD"
	ENV_PORT                string = "PORT"
	MAX_RESERVATION_TIME    int    = 5
	TOTAL_TICKETS_AVAILABLE int    = 5
)

func main() {

	// this env values will be overriden by the configs from the deployment in kubernetes
	conf := config.GetConfig(
		os.Getenv(ENV_DB_DIALECT),
		os.Getenv(ENV_DB_CONNECT_URI),
		os.Getenv(ENV_DB_USERNAME),
		os.Getenv(EVN_DB_PASSWORD))
	conf.ServerPort = os.Getenv(ENV_PORT)

	configPtr := flag.String("config", "config.local.yml", "a config file path")
	flag.Parse()

	configdata, err := ioutil.ReadFile(*configPtr)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(configdata, conf); err != nil {
		log.Fatal(err)
	}

	gormDb := database.NewGormDB(conf)
	// create guest table if not existed
	gormDb.AutoMigrate(&model.Guest{})

	guestService := registeredguest.NewGuestService(gormDb)

	// now setup/init the app controller and handlers
	logger := glog.NewLogfmtLogger(os.Stdout)
	logger = glog.With(logger, "ts", glog.DefaultTimestampUTC)
	logger = glog.With(logger, "caller", glog.DefaultCaller)
	inProgressService := inprogress.NewInMemoryService()
	inProgressService = inprogress.NewLoggingMiddlewareService(logger)(inProgressService)

	controller.Init(
		time.Duration(MAX_RESERVATION_TIME)*time.Minute,
		TOTAL_TICKETS_AVAILABLE,
		guestService,
		inProgressService,
		template.Must(template.ParseGlob("views/*")))

	var homePageHandler = controller.LoggingMiddleware(logger)(controller.DispatchHomePage)
	var reservationFormHandler = controller.LoggingMiddleware(logger)(controller.DispatchReservationForm)
	var makeReservationHandler = controller.LoggingMiddleware(logger)(controller.MakeReservation)
	var payWithCardHandler = controller.LoggingMiddleware(logger)(controller.CreateTokenAuthoringMiddleWare(controller.DoAuth)(controller.PayWithCard))

	router := mux.NewRouter()
	router.HandleFunc("/", homePageHandler).Methods("GET")
	router.HandleFunc("/reservation", reservationFormHandler).Methods("GET")
	router.HandleFunc("/reserve", makeReservationHandler).Methods("POST")
	router.HandleFunc("/charge", payWithCardHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":"+conf.ServerPort, router))
}
