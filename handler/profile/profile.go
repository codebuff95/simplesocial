package profile

import(
	"net/http"
	"log"
	"fmt"
	"simplesocial/user"
	"simplesocial/friend"
	"simplesocial/sessions"
	"html/template"
	"time"
)

type Profile struct{
	TargetUser user.User
	MyFriendship friend.Friendship
	FormSid string
	ActiveFriends *[]*string
}

func ProfileHandler(w http.ResponseWriter, r *http.Request){
	requestingUserId := user.Authenticate(r)
	if requestingUserId == ""{
		log.Println("Request Session UserId not authentic. Redirecting to login page.")
		http.Redirect(w,r, "/login",http.StatusSeeOther)
		return
	}
	log.Println("Request Session UserId is authentic. Processing profile display request.")
	targetUserId := r.URL.Path[len("/profile/"):]
	if targetUserId == requestingUserId{
		log.Println("User requesting self profile. Redirect to home page.")
		http.Redirect(w,r, "/home",http.StatusSeeOther)
		return
	}
	log.Println("Requesting profile of userid:",targetUserId)
	targetUser := user.GetUser(targetUserId)
	if targetUser == nil{
		//fmt.Fprintf(w,"User does not exist. Redirecting to home...")
		http.Redirect(w,r, "/home",http.StatusSeeOther)
		return
	}
	targetUserProfile := Profile{TargetUser : *targetUser}
	log.Println("Requesting friendship of userid1 =",requestingUserId," userid2 =",string(targetUserId))
	myFriendship := friend.GetFriendship(requestingUserId,string(targetUserId))
	if myFriendship == nil{
		log.Println("Friendship does not exist between the users.")
		targetUserProfile.MyFriendship.Status = friend.DNE
	}else{
		log.Println("Friendship exists between the users, with status = ",myFriendship.Status)
		targetUserProfile.MyFriendship = *myFriendship
	}
	targetUserProfile.ActiveFriends = friend.GetActiveFriends(targetUserId)
	t,err := template.ParseFiles("simplesocialtmp/profile.html")
	if err != nil{
		log.Println("Problem parsing profile.html")
		log.Println(err)
	}else{
		//Crating new formsid for addfriendship/deletefriendship form on profile
		thisSession := sessions.GlobalSM["formsm"].SetSession("0",time.Minute * 5) // Form valid for 5 minutes.
		if thisSession == nil{
			fmt.Fprintf(w,"Error creating session for form. Please try again in some time.")
			return
		}
		if thisSession.Status == sessions.ACTIVE{
			targetUserProfile.FormSid = thisSession.Sid
			log.Println("Generating new form to client",r.RemoteAddr,"with formsid =",thisSession.Sid)
		}
		log.Println("Successfully parsed profile.html. Executing template with Profile data",fmt.Sprintf("%+v",targetUserProfile))
		t.Execute(w,targetUserProfile)
	}
}
