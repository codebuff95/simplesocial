package login
import(
	"net/http"
	"html/template"
	"simplesocial/sessions"
	"simplesocial/form"
	"simplesocial/user"
	"log"
	"time"
	"fmt"
)

func DisplayLoginPage(w http.ResponseWriter, r *http.Request){
	thisSession := sessions.GlobalSM["formsm"].SetSession("0",time.Minute * 5) // Form valid for 5 minutes.
	if thisSession == nil{
		fmt.Fprintf(w,"Error showing login page. Please try again in some time.")
	}
	if thisSession.Status == sessions.ACTIVE{
		t,_ := template.ParseFiles("simplesocialtmp/login.html")
		t.Execute(w,thisSession.Sid)
		log.Println("Generated new form to client",r.RemoteAddr,"with formsid =",thisSession.Sid)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request){
	userIsAuthentic := user.Authenticate(r)
	if userIsAuthentic{
		//redirect to homepage.
		log.Println("Request user session is authentic. Redirecting to homepage.")
		http.Redirect(w,r, "/home",http.StatusSeeOther)
		return
	}
	//usersid not authentic. Proceed to processing request.
	if r.Method == "GET"{
			DisplayLoginPage(w,r)
	}else{ // Method == POST
		r.ParseForm()
		formIsAuthentic := form.Authenticate(r)
		if formIsAuthentic{
			//Authenticate User Begin.
			userSession := user.AuthenticateLoginAttempt(r)
			if userSession == nil || userSession.Status != sessions.ACTIVE{ //Invalid login attempt.
				//fmt.Fprintf(w,"Invalid Login attempt. Retry.\n")
				DisplayLoginPage(w,r)
			}else{ //Authentic user login. Set usersid Cookie at client.
				userSidCookie := &http.Cookie{Name : "usersid", Value : userSession.Sid}
				expiry := time.Now().Add(time.Hour*24*2) // Cookie is valid till 2 days.
				if r.Form.Get("rememberme") != ""{ //rememberme selected in login form.
					log.Println("Remember me WAS selected in form.")
					userSidCookie.Expires = expiry
				}else{
					log.Println("Remember me WAS NOT selected in form.")
				}
				http.SetCookie(w,userSidCookie)
				log.Println("usersid Cookie successfully set on client. Redirecting to homepage.")
				http.Redirect(w,r, "/home",http.StatusSeeOther)
				return
			}
			//Authenticate User End.
		}else{ //FormSid not valid.
			DisplayLoginPage(w,r)
		}
	}
}
