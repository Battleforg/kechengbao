package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
)

type Database struct {
	sync.Mutex
	Name     string
	Table    string
	Backend  string
	UserName string
	Memory   map[string]interface{}
}

type Fetch func(map[string]interface{}, *sql.Rows)

func (d *Database) Connect() *sql.DB {
	db, err := sql.Open(d.Backend, "./user.db")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return db
}

func (d *Database) InitTable() {
	db := d.Connect()
	defer db.Close()
	sqlStmt := `
	create table courseinfo (cid text, qno integer);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Println(err)
	}

	sqlStmt = `
		create table userinfo (userid text, password text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Println(err)
	}
}

func NewDB(backend, dbname, user, table string) *Database {
	ret := new(Database)
	ret.Backend = backend
	ret.Name = dbname
	ret.UserName = user
	ret.Table = table
	ret.Memory = make(map[string]interface{})
	return ret
}

func (ud *Database) Cache(f Fetch) {
	db := ud.Connect()
	defer db.Close()

	rows, err := db.Query("select * from " + ud.Table)
	if err != nil {
		log.Panic(err.Error())
	}

	for rows.Next() {
		f(ud.Memory, rows)
	}
}

func FetchCourseBrief(s map[string]interface{}, rows *sql.Rows) {
	var cinfo CourseBrief
	err := rows.Scan(&cinfo.CourseID, &cinfo.QNo)
	if err != nil {
		log.Panic(err)
	}
	s[cinfo.CourseID] = cinfo
}
