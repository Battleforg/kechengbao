package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Profile struct {
	//in-system no.
	UserID string
	//ids of courses
	Courses   []string
	Id        Ident
	Education Degree
	Info      Personal
}

type Ident struct {
	//student no.
	WorkID string
	//id card no.
	CardID string
	Email  string
	Phone  string
}

type Degree struct {
	University string
	School     string
	Major      string
	Level      string
	StartYear  string
	EndYear    string
}

type Personal struct {
	NickName  string
	Name      string
	Education []Degree
}

//delete profile by id
func DeleteProfile(cid string) {

}

//fetch profile by id
func FetchProfile(uid string) (string, error) {
	res, err := fetchProfile(uid)
	if err != nil {
		return "", err
	}
	b, _ := json.Marshal(res)
	return string(b), nil
}

func fetchProfile(uid string) (Profile, error) {
	var res Profile
	session, err := mgo.Dial(COURSE_HOST)
	if err != nil {
		return res, err
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DATABASE).C(PROFILE)

	if err != nil {
		return res, err
	}
	err = c.Find(bson.M{"userid": uid}).One(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func FetchProfileBatch(batch string) (string, error) {
	var uids []string
	err := json.Unmarshal([]byte(batch), &uids)
	if err != nil {
		return "", err
	}
	var res []Profile
	for _, v := range uids {
		p, err := fetchProfile(v)
		if err != nil {
			return "", err
		}
		res = append(res, p)
	}
	b, _ := json.Marshal(res)
	return string(b), nil
}

func CreateProfile(v Profile) error {
	session, err := mgo.Dial(COURSE_HOST)
	if err != nil {
		return err
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DATABASE).C(PROFILE)
	err = c.Remove(bson.M{"userid": v.UserID})
	if err != nil {
		log.Println(err.Error())
	}
	return c.Insert(v)
}

func AppendCourse(uid, cid string) error {
	session, err := mgo.Dial(COURSE_HOST)
	if err != nil {
		return err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DATABASE).C(PROFILE)
	err = c.Update(bson.M{"userid": uid}, bson.M{"$push": bson.M{"courses": cid}})
	return err
}

func RemoveCourse(uid, cid string) error {
	session, err := mgo.Dial(COURSE_HOST)
	if err != nil {
		return err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DATABASE).C(PROFILE)
	err = c.Update(bson.M{"userid": uid}, bson.M{"$pull": bson.M{"courses": cid}})
	return err
}
