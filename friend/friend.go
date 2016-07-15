package friend

import(
	"log"
	"simplesocial/databases"
	"net/http"
	"html/template"
	"simplesocial/form"
	"simplesocial/user"
	"time"
	"fmt"
)

const(
	ACTIVE int = 2
	PENDING int = 1
	DNE int = 3 // Does not exist.
)

type Friendship struct{
	Userid1,Userid2 string
	Status int
	Statussince string
}

// GetFriendship takes as parameters both userids, and returns the Friendship struct value between them.
func GetFriendship(userid1, userid2 string) *Friendship{
	myFriendship := Friendship{}
	log.Println("Getting requested Friendship between",userid1,"and",userid2)
	row := databases.GlobalDBM["mydb"].Con.QueryRow("SELECT statussince,CONVERT(status,CHAR(11))  FROM friendship WHERE (userid1='"+userid1+"' AND userid2='"+userid2+"') OR (userid1='"+userid2+"' AND userid2='"+userid1+"')",)
	err := row.Scan(&myFriendship.Statussince,&myFriendship.Status)
	if err != nil{
		log.Println("Friendship does not exist.","Returning nil Friendship.")
		log.Println(err)
		return nil
	}
	log.Println("Success getting Friendship.")
	myFriendship.Userid1 = userid1
	myFriendship.Userid2 = userid2
	return &myFriendship
}

func DeleteFriendship(userid1,userid2 string) bool{
	log.Println("Deleting Friendship with userid1:",userid1,"userid2:",userid2)
	stmt, err := databases.GlobalDBM["mydb"].Con.Prepare("DELETE FROM friendship WHERE (userid1='"+userid1+"' AND userid2='"+userid2+"') OR (userid1='"+userid2+"' AND userid2='"+userid1+"')")
	if err != nil{
		log.Println(err)
		return false
	}
	res,err := stmt.Exec()
	rowsdeleted, _ := res.RowsAffected()
	log.Println("deleted",rowsdeleted,"rows for friendship")
	return true
}

func AddFriendship(userid1,userid2 string) bool{
	log.Println("Adding Friendship with userid1:",userid1,"userid2:",userid2)
	addtime := time.Now().Format("2006-01-02 15:04:05")
	stmt, err := databases.GlobalDBM["mydb"].Con.Prepare("INSERT INTO friendship SET userid1 = '"+userid1+"', userid2 = '"+userid2+"', status = '"+fmt.Sprintf("%d",PENDING)+"', statussince = '"+addtime+"'")
	if err != nil{
		log.Println(err)
		return false
	}
	_,err = stmt.Exec()
	if err != nil{
		return false
	}
	return true
}

func RemoveFriendHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST"{
		//Redirect to homepage.
		http.Redirect(w,r,"/home",http.StatusSeeOther)
		return
	}
	log.Println("Remove Friend Handler invoked with POST request.")
	r.ParseForm()
	formRid := form.Authenticate(r)
	if formRid == ""{ //FormSid not valid. Redirect to Homepage.
		http.Redirect(w,r,"/home",http.StatusSeeOther)
		return
	}
	//FormSid is valid.
	requestingUserId := user.Authenticate(r)
	if requestingUserId == ""{
		log.Println("Request Session UserId not authentic. Redirecting to login page.")
		http.Redirect(w,r, "/login",http.StatusSeeOther)
		return
	}
	log.Println("Request Session UserId is authentic. Processing Remove Friend request.")
	targetUserId := r.Form.Get("targetuserid")
	targetUser := user.GetUser(targetUserId)
	if targetUser == nil{
		//fmt.Fprintf(w,"User does not exist. Redirecting to home...")
		log.Println("Target user not valid. Redirecting to homepage.")
		http.Redirect(w,r, "/home",http.StatusSeeOther)
		return
	}
	friendshipDeleted := DeleteFriendship(requestingUserId, targetUserId)
	t,_ := template.ParseFiles("simplesocialtmp/removefriend.html")
	if t == nil{
		log.Println("Problem parsing file removefriend.html")
		return
	}
	t.Execute(w,friendshipDeleted)
}

func AddFriendHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST"{
		//Redirect to homepage.
		http.Redirect(w,r,"/home",http.StatusSeeOther)
		return
	}
	log.Println("Add Friend Handler invoked with POST request.")
	r.ParseForm()
	formRid := form.Authenticate(r)
	if formRid == ""{ //FormSid not valid. Redirect to Homepage.
		http.Redirect(w,r,"/home",http.StatusSeeOther)
		return
	}
	//FormSid is valid.
	requestingUserId := user.Authenticate(r)
	if requestingUserId == ""{
		log.Println("Request Session UserId not authentic. Redirecting to login page.")
		http.Redirect(w,r, "/login",http.StatusSeeOther)
		return
	}
	log.Println("Request Session UserId is authentic. Processing Add Friend request.")
	targetUserId := r.Form.Get("targetuserid")
	targetUser := user.GetUser(targetUserId)
	if targetUser == nil{
		//fmt.Fprintf(w,"User does not exist. Redirecting to home...")
		log.Println("Target user not valid. Redirecting to homepage.")
		http.Redirect(w,r, "/home",http.StatusSeeOther)
		return
	}
	friendshipAdded := AddFriendship(requestingUserId, targetUserId)
	t,_ := template.ParseFiles("simplesocialtmp/addfriend.html")
	if t == nil{
		log.Println("Problem parsing file addfriend.html")
		return
	}
	t.Execute(w,friendshipAdded)
}

func GetActiveFriends(targetUserId string) *[]*string{
	log.Println("Getting active friends of targetuserid:",targetUserId)
	rows,err := databases.GlobalDBM["mydb"].Con.Query("SELECT CONVERT(userid1,CHAR(11)),CONVERT(userid2,CHAR(11)) FROM friendship WHERE (userid1='"+targetUserId+"' OR userid2='"+targetUserId+"') AND status='"+fmt.Sprintf("%d",ACTIVE)+"'")
	if err != nil{
		log.Println(err)
		return nil
	}
	mySlice := make([]*string,0)
	for rows.Next(){
		var userid1, userid2,reqduserid string
		err = rows.Scan(&userid1,&userid2)
		if err != nil{
			log.Println(err)
			break
		}
		if userid1 == targetUserId{
			reqduserid = userid2
		}else{
			reqduserid = userid1
		}
		mySlice = append(mySlice, &reqduserid)
	}
	log.Println("Returning active friends slice of size",len(mySlice))
	if len(mySlice) == 0{
		return nil
	}
	return &mySlice
}
