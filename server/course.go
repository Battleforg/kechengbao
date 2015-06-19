package main

import (
	"encoding/json"
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

/*
 * Storage by mongodb
 */

//client should not handle with client-specific date or time
//http.Error(w, "ss", 200)

type Chapter struct {
	No    string
	Title string
	Intro string
	Text  string //optional
	//	Resources []File
}

type ClassHour struct {
	Day        string
	StartClass string
	EndClass   string
}

type CourseHour struct {
	StartWeek  string
	EntWeek    string
	ClassHours []ClassHour
}

type Classroom struct {
	Campus   string
	Building string
	RoomNo   string
}

type Course struct {
	CourseID  string
	TimeStamp int64  `json:"-"`
	Name      string //Software Engineering
	Code      string //Course Code
	Term      string
	Hour      CourseHour
	Room      Classroom
	//	Resources []File

	Teachers []string
	TAs      []string
	Students []string

	Chapters []Chapter
	//	Process  int invisible to client
}

//explore courses

//TODO: authentication abstraction

func CreateCourse(v Course) error {
	//todo
	//allocate id
	//multiple version management
	session, err := mgo.Dial(COURSE_HOST)
	if err != nil {
		return err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DATABASE).C(COURSE)

	if err != nil {
		return err
	}
	err = c.Remove(bson.M{"courseid": v.CourseID})
	if err != nil {
		log.Println(err.Error())
	}
	err = c.Insert(v)
	if err != nil {
		return err
	}
	return nil
}

func AppendAs(cid, uid, role string) error {
	if CheckIsStudent(cid, uid) {
		return errors.New("exists")
	}
	session, err := mgo.Dial(COURSE_HOST)
	if err != nil {
		return err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DATABASE).C(COURSE)
	err = c.Update(bson.M{"courseid": cid}, bson.M{"$push": bson.M{role: uid}})
	return err
}

func RemoveUser(cid, uid string) error {
	if !CheckIsStudent(cid, uid) {
		return errors.New("not exists")
	}
	session, err := mgo.Dial(COURSE_HOST)
	if err != nil {
		return err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DATABASE).C(COURSE)
	err = c.Update(bson.M{"courseid": cid}, bson.M{"$pull": bson.M{"students": uid}})
	return err
}

//fetch course by id
func FetchCourse(cid string) (string, error) {
	var res Course
	session, err := mgo.Dial(COURSE_HOST)
	if err != nil {
		return "", err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DATABASE).C(COURSE)

	if err != nil {
		return "", err
	}
	err = c.Find(bson.M{"courseid": cid}).One(&res)
	if err != nil {
		return "", err
	}
	b, _ := json.Marshal(res)
	return string(b), nil
}

func CleanCourse() {
	session, err := mgo.Dial(COURSE_HOST)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DATABASE).C(COURSE)

	info, err := c.RemoveAll(bson.M{})
	if err != nil {
		log.Panic(err)
	}
	log.Println(info.Updated)
	log.Println(info.Removed)
}

func CleanProfile() {
	session, err := mgo.Dial(COURSE_HOST)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DATABASE).C(COURSE)

	info, err := c.RemoveAll(bson.M{})
	if err != nil {
		log.Panic(err)
	}
	log.Println(info.Updated)
	log.Println(info.Removed)
}
