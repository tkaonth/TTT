package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/testing")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8080", nil)
}

func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "index.html")
}

func loginPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "login.html")
		return
	}

}

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "signup.html")
		return
	}
	username := req.FormValue("username")
	password := req.FormValue("password")
	var user string
	fmt.Println("username : " + username)
	fmt.Println("pasword : " + password)
	fmt.Println("user : " + user)
	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)
	fmt.Println("error : " + err.Error())
	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error,unable to create your accout", 500)
			return
		}
		_, err = db.Exec("INSERT INTO users(username,password) VALUES (?,?)", username, hashedPassword)
		if err != nil {
			http.Error(res, "Server error,unable to create your accout", 500)
			return
		}
		res.Write([]byte("User created!"))
		return
	case err != nil:
		http.Error(res, "Server error,unable to create your accout", 500)
		return
	default:
		http.Redirect(res, req, "/", 301)
	}
}
