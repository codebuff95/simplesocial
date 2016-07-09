package form
import "net/http"
import "simplesocial/sessions"
import "log"
func Authenticate(r *http.Request) bool{ //NOTE: r should be parsed for forms.
	log.Println("Authenticating Form SID")
	FormSid := r.Form.Get("formsid") //Authenticate form by its formsid field.
	formIsAuthentic := sessions.GlobalSM["formsm"].Authenticate(FormSid)
	if formIsAuthentic{
		sessions.GlobalSM["formsm"].DeleteSession(FormSid)
	}
	return formIsAuthentic
}
