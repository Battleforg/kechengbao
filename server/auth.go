package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

//TODO: refactor

type User struct {
	UserID   string
	Password string
	//Session  string
	//Expire   int64
	//Role     string //Admin, User(TA, student...)
}

type UserDB struct {
	sync.Mutex
	DB Database

	//From UserID to User struct
	Users map[string]User
}

func (u *User) FromJSON(js []byte) bool {
	err := json.Unmarshal(js, u)
	if err != nil {
		log.Printf("User Unmarshal failed: %s\n", js)
		return false
	}
	return true
}

func NewUserDB(backend, dbname, user, table string) *UserDB {
	ud := new(UserDB)
	ud.DB.Backend = backend
	ud.DB.Name = dbname
	ud.DB.UserName = user
	ud.DB.Table = table

	return ud
}

func (ud *UserDB) Cache() {
	db := ud.DB.Connect()
	defer db.Close()

	rows, err := db.Query("select * from userinfo")
	if err != nil {
		log.Fatal(err.Error())
	}

	for rows.Next() {
		var uinfo User
		err = rows.Scan(&uinfo.UserID, &uinfo.Password)
		ud.Users[uinfo.UserID] = uinfo
	}
}

//Create Table
func (ud *UserDB) Init() {
	//init map
	ud.Users = make(map[string]User)
}

func (ud *UserDB) Insert(userid, password string) {
	db := ud.DB.Connect()
	defer db.Close()

	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s values ($1, $2)", ud.DB.Table))
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = stmt.Exec(userid, password)
	if err != nil {
		//don't exit
		log.Print(err.Error())
	}
	ud.Users[userid] = User{userid, password}
}

func (ud *UserDB) Clean() {
	db := ud.DB.Connect()

	stmt, err := db.Prepare(fmt.Sprintf("DELETE FROM %s", ud.DB.Table))
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = stmt.Exec()
	if err != nil {
		//don't exit
		log.Print(err.Error())
	}
}

func (ud *UserDB) CheckUserExist(uid string) bool {
	if ud.Users[uid].UserID != "" {
		return true
	} else {
		return false
	}
}

func (ud *UserDB) CheckUserAuth(uid, passwd string) bool {
	if ud.Users[uid].UserID != "" && db.Users[uid].Password == passwd {
		return true
	} else {
		return false
	}
}
