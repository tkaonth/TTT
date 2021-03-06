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

type Users struct {
	username string
	password string
}

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
	var user Users
	user = Users{
		username: req.FormValue("username"),
		password: req.FormValue("password"),
	}
	fmt.Println("username : " + user.username)
	fmt.Println("password : " + user.password)
	var id string
	var dbName string
	var dbPass string

	err = db.QueryRow("SELECT * FROM users WHERE username=?", user.username).Scan(&id, &dbName, &dbPass)
	// for selDB.Next() {
	// 	var id int
	// 	var username, password string
	// 	err = selDB.Scan(&id, &username, &password)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	fmt.Println("username : " + username)
	// 	fmt.Println("password : " + password)

	// }
	fmt.Println("username : " + dbName)
	fmt.Println("password : " + dbPass)

	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(user.password))

	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	res.Write([]byte("Hello : " + dbName))

}

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "signup.html")
		return
	}
	var user Users
	user = Users{
		username: req.FormValue("username"),
		password: req.FormValue("password"),
	}
	var userdb string
	//fmt.Println("username : " + user.username)
	//fmt.Println("password : " + user.password)
	err := db.QueryRow("SELECT username FROM users WHERE username=?", user.username).Scan(&userdb)
	//fmt.Println("userdb : " + userdb)
	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error,unable to create your accout", 500)
			return
		}
		_, err = db.Exec("INSERT INTO users(username,password) VALUES (?,?)", user.username, hashedPassword)
		if err != nil {
			http.Error(res, "Server error,unable to create your accout", 500)
			return
		}
		res.Write([]byte("User created!"))
		//http.Redirect(res, req, "/", 200)
		return
	case err == nil:
		http.Error(res, "Server error,unable to create your accout", 500)
		return
	default:
		http.Redirect(res, req, "/", 301)
		return
	}
}
