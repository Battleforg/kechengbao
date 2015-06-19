package main

/*
 * Main server
 */

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	ROOT_DIR   = "/home/kassian/pub/"
	MAX_MEMORY = 2 << 20 * 4
)

var (
	db   = NewUserDB("postgres", "course", "courseadmin", "userinfo")
	cbdb = NewCBDB("postgres", "course", "courseadmin", "courseinfo")
)

type MainServer struct {
}

func deploy(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(404)
		return
	}
	err := r.ParseMultipartForm(MAX_MEMORY)
	if err != nil {
		w.Write([]byte("MEMORY LIMIT EXCEEDED\n"))
		return
	}
	m := r.MultipartForm
	files := m.File["uploadfile"]
	for i, _ := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			w.Write([]byte("OPEN FILE FAILED\n"))
			return
		}
		log.Printf("Uploading %s from client\n", files[i].Filename)
		f, err := os.OpenFile(ROOT_DIR+files[i].Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		HandleInfo(err)
		defer f.Close()
		io.Copy(f, file)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(ROOT_DIR + "fake.html")
	HandleInfo(err)

	tmpl.Execute(w, nil)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(404)
		return
	}
	log.Println("New connection from ", r.URL)
	r.ParseForm()

	uid := r.FormValue("UserID")
	passwd := r.FormValue("Password")
	action := r.FormValue("Action")
	if uid == "" || passwd == "" || action == "" {
		w.WriteHeader(FORM_INVALID)
		return
	}

	cid := r.FormValue("CourseID")

	if action == "REGISTER" {
		if db.CheckUserExist(uid) {
			w.WriteHeader(USER_EXISTED)
		} else {
			var t Profile
			t.UserID = uid
			CreateProfile(t)
			db.Insert(uid, passwd)
			w.WriteHeader(SUCCESS)
		}
		return
	}
	if !db.CheckUserAuth(uid, passwd) {
		w.WriteHeader(AUTH_FAILED)
		return
	}

	//Now we've checked the identification
	log.Printf("Handling %s Request\n", action)
	switch action {
	case "LOGIN":
		w.WriteHeader(SUCCESS)
	case "LOGOFF":
		w.WriteHeader(SUCCESS)
	case "UPDATE_PROFILE":
		var t Profile
		json.Unmarshal([]byte(r.FormValue("Profile")), &t)
		err := CreateProfile(t)
		HTTPHandleError(w, err, nil)
	case "PROFILE":
		p, err := FetchProfile(uid)
		HTTPHandleError(w, err, []byte(p))
	case "PROFILE_BATCH":
		p, err := FetchProfileBatch(r.FormValue("UIDList"))
		HTTPHandleError(w, err, []byte(p))
	case "LEAVE":
		err := RemoveUser(cid, uid)
		if err == nil {
			err = RemoveCourse(uid, cid)
		}
		HTTPHandleError(w, err, nil)
	case "JOIN":
		err := AppendAs(cid, uid, "students")
		if err == nil {
			err = AppendCourse(uid, cid)
		}
		HTTPHandleError(w, err, nil)
	case "ACCEPT":
	case "KICK":
	case "POLL":
		t := r.FormValue("Type")

		if t == "" {
			return
		}
		p, err := RetrieveMessage(cid + "-" + t)
		HTTPHandleError(w, err, p)
	case "POST":
		var n Notif
		n.Date = time.Now().Local().String()
		n.Text = r.FormValue("Body")
		b, err := json.Marshal(n)
		if err != nil {
			return
		}

		t := r.FormValue("Type")
		if t == "" {
			return
		}
		err = AppendMessage(cid+"-"+t, string(b))
		HTTPHandleError(w, err, nil)
	case "CREATE":
		var t Course
		json.Unmarshal([]byte(r.FormValue("CourseInfo")), &t)
		err := CreateCourse(t)
		if err != nil {
			log.Print(err.Error())
			w.WriteHeader(FAILED)
			return
		}
		w.WriteHeader(SUCCESS)
		cbdb.InsertNewCourseMeta(t.CourseID)
	case "DETAIL":
		p, err := FetchCourse(cid)
		HTTPHandleError(w, err, []byte(p))
	case "LISTQ":
		p, err := FetchQuestion(uid, cid)
		log.Println(p)
		HTTPHandleError(w, err, []byte(p))
	case "ASK":
		var q Question
		q.CourseID = cid
		q.Ask = Message{uid, time.Now().Local().String(), r.FormValue("Body")}
		q.QuestionNo = strconv.Itoa(cbdb.GenerateQNo(cid))
		err := AddQuestion(uid, cid, q)
		HTTPHandleError(w, err, nil)
	case "ANWSER":
		var q Question
		q.CourseID = cid
		q.QuestionNo = r.FormValue("QuestionNo")
		q.Anwser = []Message{{uid, time.Now().Local().String(),
			r.FormValue("Body")}}

		err := AddAnwser(uid, cid, q)
		HTTPHandleError(w, err, nil)
	default:
		w.WriteHeader(UNDEFINED_ACTION)
	}
}

func main() {
	log.Println(time.Now().Local().String())

	db.Init()
	db.Cache()
	cbdb.Init()
	cbdb.Cache()
	log.Println(db.Users)
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/", sayhelloName)
	http.HandleFunc("/deploy", deploy)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(ROOT_DIR))))

	err := http.ListenAndServe(":8080", nil)
	HandleFatal(err)
}
