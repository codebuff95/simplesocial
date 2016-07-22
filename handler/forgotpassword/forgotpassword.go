package forgotpassword

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"simplesocial/form"
	"simplesocial/sessions"
	"simplesocial/user"
	"time"
)

const (
	RANDPASSLEN int = 8
)

func DisplayForgotPasswordForm(w http.ResponseWriter, r *http.Request) {
	thisSession := sessions.GlobalSM["formsm"].SetSession("0", time.Minute*5) // Form valid for 5 minutes.
	if thisSession == nil {
		fmt.Fprintf(w, "Error showing forgotpassword form. Please try again in some time.")
	}
	if thisSession.Status == sessions.ACTIVE {
		t, err := template.ParseFiles("simplesocialtmp/forgotpassword.html")
		if err != nil {
			log.Println("Problem parsing forgotpassword template.")
			return
		}
		t.Execute(w, thisSession.Sid)
		log.Println("Generated new forgotpassword form to client", r.RemoteAddr, "with formsid =", thisSession.Sid)
	}
}

func CreateNewRandomPassword() string {
	newPass := sessions.GenerateUniqueSid()
	return newPass[:RANDPASSLEN]
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	userIsAuthentic := user.Authenticate(r)
	if userIsAuthentic != "" {
		//redirect to homepage.
		log.Println("Request user session is authentic. Redirecting to homepage.")
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	//usersid not authentic. Proceed to processing request.
	if r.Method == "GET" {
		DisplayForgotPasswordForm(w, r)
		return
	}
	//Method = POST.
	r.ParseForm()
	formIsAuthentic := form.Authenticate(r)
	if formIsAuthentic == "" { //submitted form is not authentic.
		DisplayForgotPasswordForm(w, r)
		return
	}
	requestedUser := user.AuthenticateForgotPasswordAttempt(r)
	if requestedUser == nil {
		log.Println("User with entered email does not exist.")
	} else {
		log.Println("User with entered email exists.")
		newPassword := user.ResetPassword(requestedUser, CreateNewRandomPassword())
		if newPassword == "" {
			log.Println("User password not changed.")
		} else {
			log.Println("Successfully changed user password.")

			err := user.ResetPasswordEmail(requestedUser, newPassword)
			if err != nil {
				log.Println("Could not send Reset Password Email:", err)
			} else {
				log.Println("Successfully sent Reset Password Email")
			}
		}
	}
}
