package database

import (
	"database/sql"

	"github.com/Benchkram/errz"
	"github.com/ydsxiong/go-playground/people"
)

// FindAllUsers fetch a list of all users existing in db
func FindAllUsers(db *sql.DB) (allUsers []*people.User, err *error) {
	allUsers = []*people.User{}
	defer errz.Recover(err)

	// query
	rows, e := db.Query("SELECT * FROM users")
	err = &e
	errz.Fatal(*err)
	defer rows.Close()

	for rows.Next() {
		var user people.User
		//var name, dep, created sql.NullString
		//e = rows.Scan(&user.Uid, &name, &dep, &created)
		err = &e
		errz.Fatal(*err)

		user.Username = "" //name.String
		// user.Department = dep.String
		// user.Created = created.String
		allUsers = append(allUsers, &user)
	}
	return allUsers, err
}

// FindUserById fetch a user by its ID
func FindUserById(uid string) (user *people.User, err *error) {
	defer errz.Recover(err)

	_, err = FetchDBFromCache()
	errz.Fatal(*err)

	user = &people.User{}
	// var name, dep, created sql.NullString
	// // query
	// row := db.QueryRow("SELECT * FROM users WHERE uid=?", uid)
	// e := row.Scan(&user.Uid, &name, &dep, &created)

	// if *err != nil && *err != sql.ErrNoRows {
	// 	err = &e
	// 	errz.Fatal(*err)
	// }
	// user.Username = name.String
	// user.Department = dep.String
	// user.Created = created.String
	return user, err
}

// AddNewUser create and add a new user into DB
func AddNewUser(name string, department string) (err *error) {
	defer errz.Recover(err)

	db, err := FetchDBFromCache()
	errz.Fatal(*err)

	stmt, e := db.Prepare("INSERT INTO users (username, department) VALUES(?,?)")
	if e != nil {
		err = &e
		errz.Fatal(*err)
	}
	stmt.Exec(name, department)

	return err
}

// UpdateUser update user's data into DB
func UpdateUser(uid string, name string, department string) (err *error) {
	defer errz.Recover(err)

	db, err := FetchDBFromCache()
	errz.Fatal(*err)

	stmt, e := db.Prepare("update users set username=?, department=? where uid=?")
	if e != nil {
		err = &e
		errz.Fatal(*err)
	}
	stmt.Exec(name, department, uid)

	return err
}

// DeleteUser delete an user from the DB
func DeleteUser(uid string) (err *error) {
	defer errz.Recover(err)

	db, err := FetchDBFromCache()
	errz.Fatal(*err)

	stmt, e := db.Prepare("DELETE FROM users where uid=?")
	if e != nil {
		err = &e
		errz.Fatal(*err)
	}
	stmt.Exec(uid)

	return err
}
