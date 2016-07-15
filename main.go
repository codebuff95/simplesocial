package main
import(
	"net/http"
	"log"
	"math/rand"
	"time"
	//"html/template"
	//"fmt"
	"simplesocial/handler/login"
	"simplesocial/handler/home"
	"simplesocial/handler/profile"
	"simplesocial/handler/register"
	"simplesocial/friend"
	"simplesocial/user"
	"simplesocial/databases"
	"simplesocial/sessions"
)
func main(){
	var err error
	rand.Seed(time.Now().UTC().UnixNano())
	databases.InitGlobalDBM()
	sessions.InitGlobalSM()
	databases.GlobalDBM["mydb"] = &databases.DBManager{Name : "mysql", Database: "test", User : "root", Password : "123456"}
	err = databases.GlobalDBM["mydb"].Open()
	if err != nil{
		log.Fatal(err)
	}else{
	log.Println("Successfully Opened mydb!")
	}
	//Initialise User Session Manager in sessions.(GLobal Session Managers Map) using mydb Database and HARDCODED Tablename.
	sessions.GlobalSM["usersm"] = &sessions.SessionManager{Db : databases.GlobalDBM["mydb"], TableName : "usersession"}
	sessions.GlobalSM["formsm"] = &sessions.SessionManager{Db : databases.GlobalDBM["mydb"], TableName : "formsession"}
	http.HandleFunc("/home",home.HomeHandler)
	http.HandleFunc("/login",login.LoginHandler)
	http.HandleFunc("/logout",user.LogoutHandler)
	http.HandleFunc("/register",register.RegisterHandler)
	http.HandleFunc("/profile/",profile.ProfileHandler)
	http.HandleFunc("/removefriend",friend.RemoveFriendHandler)
	http.HandleFunc("/acceptfriend",friend.AcceptFriendHandler)
	http.HandleFunc("/addfriend",friend.AddFriendHandler)
	http.ListenAndServe(":8080",nil)
}
