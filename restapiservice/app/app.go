package app

import (
	"fmt"
	"net/http"

	"github.com/Benchkram/errz"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/ydsxiong/go-playground/restapiservice/app/config"
	"github.com/ydsxiong/go-playground/restapiservice/app/handler"
	"github.com/ydsxiong/go-playground/restapiservice/app/model"
)

// App has router and db instances, a global variable to hold on to the db connection and router
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// App initialize with predefined configuration
func (a *App) Initialize(config *config.Config) {
	pwd := config.DB.Password
	if pwd != "" {
		pwd = ":" + pwd
	}
	dbURI := fmt.Sprintf("%s%s@%s",
		config.DB.Username,
		pwd,
		config.DB.ConnectUri)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	errz.Fatal(err, "Could not connect database\n")
	defer errz.Recover(&err)
	// if err != nil {
	// 	log.Fatal("Could not connect database")
	// }

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()
}

// Set all required routers
func (a *App) setRouters() {
	// Routing for handling the projects
	a.Get("/employees", a.GetAllEmployees)
	a.Post("/employees", a.CreateEmployee)
	a.Get("/employees/{name}", a.GetEmployee)
	a.Put("/employees/{name}", a.UpdateEmployee)
	a.Delete("/employees/{name}", a.DeleteEmployee)
	a.Put("/employees/{name}/disable", a.DisableEmployee)
	a.Put("/employees/{name}/enable", a.EnableEmployee)
}

// Run the app on it's router
func (a *App) Run(httpPort string) {
	err := http.ListenAndServe(":"+httpPort, a.Router)
	errz.Fatal(err)
	errz.Recover(&err)
}

// Wrap the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Wrap the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Wrap the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Wrap the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Handlers to manage Employee Data
func (a *App) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	handler.GetAllEmployees(a.DB, w, r)
}

func (a *App) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	handler.CreateEmployee(a.DB, w, r)
}

func (a *App) GetEmployee(w http.ResponseWriter, r *http.Request) {
	handler.GetEmployee(a.DB, w, r)
}

func (a *App) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	handler.UpdateEmployee(a.DB, w, r)
}

func (a *App) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	handler.DeleteEmployee(a.DB, w, r)
}

func (a *App) DisableEmployee(w http.ResponseWriter, r *http.Request) {
	handler.DisableEmployee(a.DB, w, r)
}

func (a *App) EnableEmployee(w http.ResponseWriter, r *http.Request) {
	handler.EnableEmployee(a.DB, w, r)
}
