package profile

import(
	"net/http"
	"log"
	"simplesocial/user"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request){
	userIsAuthentic := user.Authenticate(r)
	if !userIsAuthentic{
		log.Println("Request Session UserId not authentic. Redirecting to login page.")
		http.Redirect(w,r, "/login",http.StatusSeeOther)
		return
	}
	log.Println("Request Session UserId is authentic. Processing profile display request.")
	targetUserId := r.URL.Path[len("/profile/"):]
	log.Println("Requesting profile of userid:",targetUserId)
	targetUser := user.GetUser(targetUserId)
	if targetUser == nil{
		//fmt.Fprintf(w,"User does not exist. Redirecting to home...")
		http.Redirect(w,r, "/home",http.StatusSeeOther)
	}
}
