package main

import (
	"log"
	"math/rand"
	"net/http"
	"simplesocial/databases"
	"simplesocial/email"
	"simplesocial/friend"
	"simplesocial/handler/forgotpassword"
	"simplesocial/handler/home"
	"simplesocial/handler/login"
	"simplesocial/handler/profile"
	"simplesocial/handler/register"
	"simplesocial/handler/welcome"
	"simplesocial/sessions"
	"simplesocial/user"
	"time"
)

func main() {
	var err error
	rand.Seed(time.Now().UTC().UnixNano())
	databases.InitGlobalDBM()
	sessions.InitGlobalSM()
	databases.GlobalDBM["mydb"] = &databases.DBManager{Name: "mysql", Database: "test", User: "root", Password: "tanmay123"}
	err = databases.GlobalDBM["mydb"].Open()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Opened mydb!")
	}
	//Seperate goroutins for each Table Cleaner.
	go databases.GlobalDBM["mydb"].CleanTable("formsession")
	go databases.GlobalDBM["mydb"].CleanTable("usersession")
	//Initialise User Session Manager in sessions.(GLobal Session Managers Map) using mydb Database and HARDCODED Tablename.
	sessions.GlobalSM["usersm"] = &sessions.SessionManager{Db: databases.GlobalDBM["mydb"], TableName: "usersession"}
	sessions.GlobalSM["formsm"] = &sessions.SessionManager{Db: databases.GlobalDBM["mydb"], TableName: "formsession"}

	//Initialise Global Email Manager.
	email.InitGlobalEM()
	//Initialise Email Templates (WelcomeEmailTemplate,ResetPasswordTemplate).
	err = email.InitEmailTemplates()
	if err != nil {
		log.Fatal("Could not initialise Welcome Email Template.")
	}
	/*
		http.HandleFunc("/welcome", welcome.WelcomeHandler)
		http.HandleFunc("/home", home.HomeHandler)
		http.HandleFunc("/login", login.LoginHandler)
		http.HandleFunc("/logout", user.LogoutHandler)
		http.HandleFunc("/register", register.RegisterHandler)
		http.HandleFunc("/profile/", profile.ProfileHandler)
		http.HandleFunc("/removefriend", friend.RemoveFriendHandler)
		http.HandleFunc("/acceptfriend", friend.AcceptFriendHandler)
		http.HandleFunc("/addfriend", friend.AddFriendHandler)
	*/
	http.HandleFunc("/", MyHandler)
	http.HandleFunc("/profile/", profile.ProfileHandler)
	http.ListenAndServe(":8080", nil)
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("***Entering MyHandler***")
	requestPath := r.URL.Path
	if requestPath == "/home" {
		log.Println("###Home Handler###")
		home.HomeHandler(w, r)
		return
	}
	if requestPath == "/login" {
		log.Println("###Login Handler###")
		login.LoginHandler(w, r)
		return
	}
	if requestPath == "/logout" {
		log.Println("###Logout Handler###")
		user.LogoutHandler(w, r)
		return
	}
	if requestPath == "/register" {
		log.Println("###Register Handler###")
		register.RegisterHandler(w, r)
		return
	}
	/*if requestPath == "/profile/" {
		log.Println("###Profile Handler###")
		profile.ProfileHandler(w, r)
		return
	}*/
	if requestPath == "/removefriend" {
		log.Println("###RemoveFriend Handler###")
		friend.RemoveFriendHandler(w, r)
		return
	}
	if requestPath == "/acceptfriend" {
		log.Println("###AcceptFriend Handler###")
		friend.AcceptFriendHandler(w, r)
		return
	}
	if requestPath == "/addfriend" {
		log.Println("###AddFriend Handler###")
		friend.AddFriendHandler(w, r)
		return
	}
	if requestPath == "/forgotpassword" {
		log.Println("###Forgot Password Handler###")
		forgotpassword.ForgotPasswordHandler(w, r)
		return
	}
	log.Println("Request with path", r.URL.Path, "is triggering welcome handler.")
	log.Println("###Welcome Handler###")
	welcome.WelcomeHandler(w, r)
	return
}
