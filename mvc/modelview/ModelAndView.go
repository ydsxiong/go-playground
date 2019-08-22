package modelview

import (
	"log"
	"net/http"
	"text/template"

	"github.com/Benchkram/errz"
	"github.com/ydsxiong/go-playground/database"
)

var tmpl = template.Must(template.ParseGlob("views/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db, err := database.FetchDBFromCache()
	if *err == nil {
		dbUsers, err := database.FindAllUsers(db)
		if *err == nil {
			tmpl.ExecuteTemplate(w, "Index", dbUsers)
		}
	}
	errz.Log(*err, "unable to fetch user from database!")
}

func Show(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("id")
	user, err := database.FindUserById(uid)
	if *err == nil {
		tmpl.ExecuteTemplate(w, "Show", user)
	}
	errz.Log(*err, "unable to find user for the given id!")
}

func Edit(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("id")
	user, err := database.FindUserById(uid)
	if *err == nil {
		tmpl.ExecuteTemplate(w, "Edit", user)
	}
	errz.Log(*err, "unable to find user for the given id!")
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		dep := r.FormValue("department")
		err := database.AddNewUser(name, dep)
		if *err == nil {
			log.Println("INSERTED: Name: " + name + " | department: " + dep)
		} else {
			errz.Log(*err, "Failed to add the new user into system!")
		}
	}
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		department := r.FormValue("department")
		id := r.FormValue("id")
		err := database.UpdateUser(id, name, department)
		if *err == nil {
			log.Println("UPDATED: Name: " + name + " | department: " + department)
		} else {
			errz.Log(*err, "Failed to update the user data into system!")
		}
	}
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("id")
	err := database.DeleteUser(uid)
	if *err == nil {
		log.Println("DELETED")
	} else {
		errz.Log(*err, "Failed to delete the user from system!")
	}
	http.Redirect(w, r, "/", 301)
}
