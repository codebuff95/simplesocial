package user

import(
	"net/http"
	"log"
	"simplesocial/sessions"
	"simplesocial/databases"
	"time"
	"database/sql"
)
type User struct{
	Userid string
	Firstname string
	Lastname string
	Email string
	Photoid string
	Registeredon string
}
func Authenticate(r *http.Request) bool{ //Authenticate UserSID field of received request. NOTE: r should be parsed for forms.
	log.Println("Authenticating User SID on received request.")
	userCookie,err := r.Cookie("usersid")
	if err != nil{ //usersid cookie not set on client.
		log.Println("User SID authentication failed")
		return false
	}
	userSid := userCookie.Value //Authenticate user by its usersid field.
	return sessions.GlobalSM["usersm"].Authenticate(userSid)
}
func AuthenticateLoginAttempt(r *http.Request) *sessions.Session{ // NOTE: r should be parsed for forms.
	var userid string
	log.Println("Authenticating Login credentials.")
	attemptEmail := r.Form.Get("email")
	attemptPassword := r.Form.Get("password")
	log.Println("Attempt email :",attemptEmail,"Attempt Password:",attemptPassword)
	row := databases.GlobalDBM["mydb"].Con.QueryRow("SELECT userid FROM user WHERE email = '"+attemptEmail+"' AND password = '"+attemptPassword+"'")
	err := row.Scan(&userid)
	if err != nil{ // User does not exist.
		log.Println("User authentication failed.")
		return &sessions.Session{Status: sessions.DELETED}
	}else{ //User exists.
		log.Println("User authentication successful. Creating new Session.")
		return sessions.GlobalSM["usersm"].SetSession(userid,time.Hour*24*3) // Session lives in DB for 3 days.
	}
}

func GetUser(targetuserid string) *User{
	targetUser := User{}
	log.Println("Getting requested userid",targetuserid)
	row := databases.GlobalDBM["mydb"].Con.QueryRow("SELECT CONVERT(userid,CHAR(11)),firstname,lastname,CONVERT(photoid,CHAR(11)),email,registeredon FROM user WHERE userid = '"+targetuserid+"'")
	var photoid sql.NullString //handling case when photoid is null.
	err := row.Scan(&targetUser.Userid,&targetUser.Firstname,&targetUser.Lastname,&photoid,&targetUser.Email, &targetUser.Registeredon)
	if photoid.Valid{ //if photoid is not null
		targetUser.Photoid = photoid.String
	}
	if err != nil{
		log.Println("Could not find target userid",targetuserid,", Returning nil User.")
		log.Println(err)
		return nil
	}
	log.Println("Success getting target userid",targetuserid,", Returning User.")
	return &targetUser
}
