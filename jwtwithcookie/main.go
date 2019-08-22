package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/ydsxiong/golang/database"
	"github.com/ydsxiong/golang/people"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

const (
	maxLineLength = 76
	upperhex      = "0123456789ABCDEF"
	cmdSendmail   = "/usr/sbin/sendmail"
	mailFrom      = "mail@golang.com"
	mailSubject   = "Reset password"

	forgotPasswordForm = `
	<html>
	<head>
	<title>Login</title>
	</head>
	<body>
	{{if .msg}}
	<font size="3" color="red">{{.msg}}</font>
	{{end}}
	<form action="/requestpassword" name=f method="POST">
	<br><br>
	<label>Email to be sent link to: <label><input maxLength=100 size=30 name=email value="" title="email">	<br><br>
	<input type=submit value="Request link" name=email>
	</form>
	</body>
	</html>
	`

	pwdResetForm = `
	<html>
	<head>
	<title>Reset password</title>
	</head>
	<body>
	{{if .msg}}
	<font size="3" color="red">{{.msg}}</font>
	<br>
	<br>
	{{end}}
	<form action="/resetpassword" name=f method="POST">
	<label>New password: <label><input maxLength=100 size=30 name=pwd value="" title="New password">	<br><br>
	<label>Confirm password: <label><input maxLength=100 size=30 name=pwdconfirm value="" title="Confirm new password">	<br><br>
	<input type=submit value="Reset" name=reset>
	{{if .}}
	<input type="hidden" id="uuid" name="uuid" value="{{.uuid}}"/>
	{{end}}
	</form>
	</body>
	</html>
	`
	loginForm = `
	<html>
	<head>
	<title>Login</title>
	</head>
	<body>
	{{if .msg}}
	<font size="3" color="red">{{.msg}}</font>
	<br>
	<br>
	{{end}}
	<form action="/dologin" name=f method="POST">
	<label>Username: <label><input maxLength=100 size=30 name=username value="" title="username">	<br><br>
	<label>Password: <label><input maxLength=100 size=30 name=pwd value="" title="password">	<br><br>
	<input type=submit value="Sign In" name=signin><br><br>
	<a href="/forgotpwd">Forgotten password</a>
	</form>
	</body>
	</html>
	`
	homePage = `
	<html>
	<head>
	<title>Home</title>
	</head>
	<body>
	<font size="5" color="blue">Welcome: {{.user}}</font> <a href="/logout">[Logout]</a>
	<br><br>
	<p>{{.greetings}}
	<p>Feel free to browser wherever you like here!</p>
	</body>
	</html>
	`
)

var db *gorm.DB

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

var dispatchForgottenPwdForm = template.Must(template.New("reset").Parse(forgotPasswordForm))

var dispatchResetPwdForm = template.Must(template.New("reset").Parse(pwdResetForm))

var dispatchLoginForm = template.Must(template.New("login").Parse(loginForm))

var dispatchHomePage = template.Must(template.New("home").Parse(homePage))

var resettingpwds = make(map[string]string)

type Credentials struct {
	Password string `json:"password", db:"password"`
	Username string `json:"username", db:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Exception struct {
	Message string `json:"message"`
}

func main() {
	db = database.InitDBLocal()

	router := mux.NewRouter()
	router.HandleFunc("/signup", Signup).Methods("POST")
	router.HandleFunc("/login", tokenAuthoringMiddleWare(GetLoginForm)).Methods("GET")
	router.HandleFunc("/dologin", Login).Methods("POST")
	router.HandleFunc("/signin", SigninViaAPI).Methods("POST")
	router.HandleFunc("/", tokenAuthoringMiddleWare(Home)).Methods("GET")
	router.HandleFunc("/refresh", tokenAuthoringMiddleWare(Refresh)).Methods("GET")
	router.HandleFunc("/logout", tokenAuthoringMiddleWare(Logout)).Methods("GET")
	router.HandleFunc("/forgotpwd", GetForgottonPwdForm).Methods("GET")
	router.HandleFunc("/requestpassword", ForgotPassword).Methods("POST")
	router.HandleFunc("/user/{uuid:[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}}/reset", GetResetForm).Methods("GET")
	router.HandleFunc("/resetpassword", ResetPassword).Methods("POST")

	log.Fatal(http.ListenAndServe(":8008", router))
}

func tokenAuthoringMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		validateToken(w, req, next)
	})
}

func Signup(w http.ResponseWriter, r *http.Request) {
	user := &people.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil || user.Email == "" {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	existingUser := getUserByName(user.Username)
	if existingUser != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(Exception{Message: "Username already taken, please specify a different one!"})
		return
	}
	existingUser = getUserByEmail(user.Email)
	if existingUser != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(Exception{Message: "This email has already been signed up in our system!"})
		return
	}
	user.Password = hashPassword(user.Password)
	if err = db.Save(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func hashPassword(pwd string) string {
	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hased, _ := bcrypt.GenerateFromPassword([]byte(pwd), 8)
	return string(hased)
}

func GetLoginForm(w http.ResponseWriter, r *http.Request) {
	claims := retriveClaimsFromCookie(r)
	if claims != nil {
		refreshTokenIfTooCloseToExpiryTime(w, claims)
		w.Write([]byte(fmt.Sprintf("Hi %s, you are signed in already, please continue browsing your way through the site!", claims.Username)))
		return
	}
	dispatchLoginForm.Execute(w, nil)
}

func GetForgottonPwdForm(w http.ResponseWriter, r *http.Request) {
	dispatchForgottenPwdForm.Execute(w, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	handleMsg := func(msg string, code int) {
		data := map[string]string{"msg": msg}
		dispatchLoginForm.Execute(w, data)
	}

	username := r.FormValue("username")
	pwd := r.FormValue("pwd")
	if username == "" || pwd == "" {
		handleMsg("Error: Incorrect login details!", 0)
		return
	}
	doSignin(w, r, username, pwd, handleMsg)
}

func SigninViaAPI(w http.ResponseWriter, r *http.Request) {
	handleMsg := func(msg string, code int) {
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(Exception{Message: msg})
	}
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		handleMsg("Missing credentials", http.StatusBadRequest)
		return
	}

	doSignin(w, r, creds.Username, creds.Password, handleMsg)
}

func doSignin(w http.ResponseWriter, r *http.Request, username, pwd string, handleErrorMsg func(msg string, code int)) {
	storedCreds := getUserByName(username)
	if storedCreds == nil {
		storedCreds = getUserByEmail(username)
	}

	if storedCreds == nil {
		handleErrorMsg("Unrecognized user", http.StatusNotFound)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(pwd)); err != nil {
		// If the two passwords don't match, return a 401 status
		handleErrorMsg("Incorrect login details", http.StatusUnauthorized)
		return
	}

	generateToken(w, &Claims{
		Username:       storedCreds.Username,
		StandardClaims: jwt.StandardClaims{},
	})
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// to be forwarded to by the token validation handler
func Home(w http.ResponseWriter, r *http.Request) {
	claims := retriveClaimsFromCookie(r)
	if claims != nil {
		refreshTokenIfTooCloseToExpiryTime(w, claims)
		var greetings string
		if claims.Username == "test" {
			greetings = "We appreciate your custom here!"
		} else {
			greetings = "You are highly appreciated as a valuable customer here!"
		}
		data := map[string]string{"user": claims.Username, "greetings": greetings}
		dispatchHomePage.Execute(w, data)
		// w.Write([]byte(fmt.Sprintf("Welcome %s!", )))
	}
}

/**
To minimize misuse of a JWT, the expiry time is usually kept in the order of a few minutes.
Typically the client application would refresh the token in the background.
*/
// to be forwarded to by the token validation handler
func Refresh(w http.ResponseWriter, r *http.Request) {
	claims := retriveClaimsFromCookie(r)
	if claims == nil {
		return
	}
	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	refreshed, expireInseconds := refreshTokenIfTooCloseToExpiryTime(w, claims)
	if !refreshed {
		w.WriteHeader(http.StatusBadRequest)
		min, sec := expireInseconds/60, expireInseconds%60
		msg := fmt.Sprintf("Token can only be refreshed within %d seconds of being expired, remaining time: %d minutes %d seconds", 30, min, sec)
		json.NewEncoder(w).Encode(Exception{Message: msg})
		return
	}
}

// to be forwarded to by the token validation handler
func Logout(w http.ResponseWriter, r *http.Request) {
	clearoutTokenFromCookie(w)
	claims := retriveClaimsFromCookie(r)
	w.Write([]byte(fmt.Sprintf("Bye %s, see you next time when you sign in again!", claims.Username)))
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	/* from rest api client
	type emailaddr struct {
		Email string `json:"email"`
	}
	var aux = emailaddr{}
	err := json.NewDecoder(r.Body).Decode(&aux)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Exception{Message: "Missing email address"})
		return
	}
	existingUser := getUserByEmail(aux.Email)
	*/
	// this is for go template engine form:
	email := r.FormValue("email")
	if email == "" {
		data := map[string]string{"msg": "Error: enter a valid email address!"}
		dispatchForgottenPwdForm.Execute(w, data)
		return
	}

	existingUser := getUserByEmail(email)
	if existingUser == nil {
		/* for rest api response
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(Exception{Message: "This email is not recognized in our system!"})
		*/
		// for template engine form:
		data := map[string]string{"msg": "Error: this email is not recognized in our system!"}
		dispatchForgottenPwdForm.Execute(w, data)
		return
	}

	reset_uid, _ := uuid.NewV4() // 5b52d72c-82b3-4f8e-beb5-437a974842c
	link := "http://" + r.Host + "/user/" + reset_uid.String() + "/reset"

	fmt.Println(link)

	emailBody := "To reset your password, please click on the link: <a href=\"" + link +
		"\">" + link + "</a><br><br>Best Regards,<br>GoLang team"

	m := gomail.NewMessage()
	m.SetHeader("From", mailFrom)
	m.SetHeader("To", email)
	m.SetHeader("Subject", mailSubject)
	m.SetBody("text/html", emailBody)

	err := sendEMail(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Exception{Message: "cookies are required to support the request"})
	} else {
		w.Write([]byte(fmt.Sprintf("A link has been sent to this email address %s, please click on the link provided in the email to reset your password", email)))

		resettingpwds[reset_uid.String()] = email
	}
}

func GetResetForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	if _, ok := resettingpwds[uuid]; !ok {
		w.WriteHeader(http.StatusFound)
		json.NewEncoder(w).Encode(Exception{Message: "The link is no longer valid, please request a new one for resetting your password"})
		return
	}
	data := map[string]string{"uuid": uuid}
	dispatchResetPwdForm.Execute(w, data)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	newPwd := r.FormValue("pwd")
	confirmPwd := r.FormValue("pwdconfirm")
	uuid := r.FormValue("uuid")
	if newPwd == "" || newPwd != confirmPwd {
		data := map[string]string{"uuid": uuid, "msg": "Error: new passwords don't match up with each other!"}
		dispatchResetPwdForm.Execute(w, data)
		return
	}
	if userEmail, ok := resettingpwds[uuid]; ok {
		existingUser := getUserByEmail(userEmail)
		if existingUser == nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Exception{Message: "Unable to reset password, please request a new link for it"})
			return
		}

		if err := db.Model(&existingUser).UpdateColumn("password", hashPassword(newPwd)).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		delete(resettingpwds, uuid)
		w.Write([]byte("The password has been reset successfully!"))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Exception{Message: "Unable to reset password, please request a new link for it"})
	}
}

func generateToken(w http.ResponseWriter, claims *Claims) bool {
	expirationTime := time.Now().Add(3 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	return true
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "redirecting",
		Value: "true",
	})
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func validateToken(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// We can obtain the session token from the requests cookies, which come with every request
	yes, err := req.Cookie("redirecting")
	if err == nil && "true" == yes.Value {
		http.SetCookie(w, &http.Cookie{
			Name:  "redirecting",
			Value: "",
		})
		next(w, req)
		return
	}
	c, err := req.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			redirectToLogin(w, req)
			// If the cookie is not set, return an unauthorized status
			//w.WriteHeader(http.StatusUnauthorized)
			//json.NewEncoder(w).Encode(Exception{Message: "cookies are required to support the authentication"})
			return
		}
		redirectToLogin(w, req)
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			redirectToLogin(w, req)
			//w.WriteHeader(http.StatusUnauthorized)
			return
		}
		redirectToLogin(w, req)
		//w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		redirectToLogin(w, req)
		//w.WriteHeader(http.StatusUnauthorized)
		return
	}
	context.Set(req, "claimFromCookie", claims)
	next(w, req)
}

func retriveClaimsFromCookie(r *http.Request) *Claims {
	cookieclaims := context.Get(r, "claimFromCookie")
	claims, ok := cookieclaims.(*Claims)
	if ok {
		return claims
	}
	return nil
}

// check to see if the claim is to close to its expiry time,
// if yes, then refresh it, and return true to indicate so, otherwise return false to signify no change.
func refreshTokenIfTooCloseToExpiryTime(w http.ResponseWriter, claims *Claims) (bool, int) {
	var timeLeft = time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
	expireInseconds := int(timeLeft.Seconds())
	if expireInseconds > 30 { // or timeLeft > 30*time.Second
		return false, expireInseconds
	}
	return generateToken(w, claims), expireInseconds
}

func clearoutTokenFromCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now(),
	})
	w.WriteHeader(http.StatusOK)
}

func getUserByName(username string) *people.User {
	user := people.User{}
	if err := db.First(&user, people.User{Username: username}).Error; err != nil {
		return nil
	}
	return &user
}

func getUserByEmail(email string) *people.User {
	user := people.User{}
	if err := db.First(&user, people.User{Email: email}).Error; err != nil {
		return nil
	}
	return &user
}

func sendEMail(m *gomail.Message) (err error) {
	cmd := exec.Command(cmdSendmail, "-t")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	pw, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	var errs [3]error
	_, errs[0] = m.WriteTo(pw)
	errs[1] = pw.Close()
	errs[2] = cmd.Wait()
	for _, err = range errs {
		if err != nil {
			return
		}
	}
	return
}
