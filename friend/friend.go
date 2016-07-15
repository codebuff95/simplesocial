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

//AddFriendship creates a new entry in table friendship with status = PENDING.
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

//AcceptFriendship updates entry in table friendship from status = PENDING to status = ACTIVE.
func AcceptFriendship(userid1,userid2 string) bool{
	log.Println("Modifying Friendship with userid1:",userid1,"userid2:",userid2)
	addtime := time.Now().Format("2006-01-02 15:04:05")
	stmt, err := databases.GlobalDBM["mydb"].Con.Prepare("UPDATE friendship set status='"+fmt.Sprintf("%d",ACTIVE)+"',statussince='"+addtime+"' WHERE userid1='"+userid1+"' AND userid2='"+userid2+"'")
	if err != nil{
		log.Println(err)
		return false
	}
	res,err := stmt.Exec()
	if err != nil{
		return false
	}
	affect,err := res.RowsAffected()
	log.Println("Modified",affect)
	return true
}

func AcceptFriendHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST"{
		//Redirect to homepage.
		http.Redirect(w,r,"/home",http.StatusSeeOther)
		return
	}
	log.Println("Accept Friend Handler invoked with POST request.")
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
	log.Println("Request Session UserId is authentic. Processing Accept Friend request.")
	targetUserId := r.Form.Get("targetuserid")
	targetUser := user.GetUser(targetUserId)
	if targetUser == nil{
		//fmt.Fprintf(w,"User does not exist. Redirecting to home...")
		log.Println("Target user not valid. Redirecting to homepage.")
		http.Redirect(w,r, "/home",http.StatusSeeOther)
		return
	}
	friendshipAccepted := AcceptFriendship(targetUserId,requestingUserId)
	t,_ := template.ParseFiles("simplesocialtmp/acceptfriend.html")
	if t == nil{
		log.Println("Problem parsing file addfriend.html")
		return
	}
	t.Execute(w,friendshipAccepted)
}

//GetFriends returns a slice of string, each string being the userid of friend, friendship of which has a status 'status'.
//Here, the targetuserid can appear on either of userid1 or userid2 of friendship table.
func GetFriends(targetUserId string,status int) *[]*string{
	log.Println("Getting active friends of targetuserid:",targetUserId)
	rows,err := databases.GlobalDBM["mydb"].Con.Query("SELECT CONVERT(userid1,CHAR(11)),CONVERT(userid2,CHAR(11)) FROM friendship WHERE (userid1='"+targetUserId+"' OR userid2='"+targetUserId+"') AND status='"+fmt.Sprintf("%d",status)+"'")
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

//GetReceivedPendingFriends returns a slice of string, each string being the userid of friend, friendship of which has a status 'pending'.
//Here, the targetuserid can appear only on userid2 of friendship table, since we need only RECEIVED pending requests.
func GetReceivedPendingFriends(userId string) *[]*string{
	log.Println("Getting received pending friends of Userid:",userId)
	rows,err := databases.GlobalDBM["mydb"].Con.Query("SELECT CONVERT(userid1,CHAR(11)) FROM friendship WHERE userid2='"+userId+"' AND status='"+fmt.Sprintf("%d",PENDING)+"'")
	if err != nil{
		log.Println(err)
		return nil
	}
	mySlice := make([]*string,0)
	for rows.Next(){
		var reqduserid string
		err = rows.Scan(&reqduserid)
		if err != nil{
			log.Println(err)
			break
		}
		mySlice = append(mySlice, &reqduserid)
	}
	log.Println("Returning received pending friends slice of size",len(mySlice))
	if len(mySlice) == 0{
		return nil
	}
	return &mySlice
}
