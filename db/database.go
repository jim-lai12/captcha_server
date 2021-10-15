package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	_"github.com/lib/pq"
	"log"
	"time"
)



func GenerateRandomString(s int) (string, error) {
	b := make([]byte, 99)//best to use 3X it will not appear =
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:s], err
}



type DbServer struct {
	Host string `yaml:"DbHost"`
	Port string `yaml:"DbPort"`
	User string `yaml:"DbUser"`
	Password string `yaml:"DbPassword"`
	Db map[string]*sql.DB

}

/* if not use config to create
func (dbserver DbServer) Initialize(host string,port,user string,password string)  {
	dbserver.Host = host
	dbserver.Port = port
	dbserver.User = user
	dbserver.Password = password
	psql := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s sslmode=disable",
		dbserver.Host, dbserver.Port, dbserver.User, dbserver.Password)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		panic(err)
	}
	dbserver.Db["MainDbConnect"] = db
}
 */

func (dbserver *DbServer) Initialize()  {
	psql := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=postgres sslmode=disable",
		dbserver.Host, dbserver.Port, dbserver.User, dbserver.Password)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		panic(err)
	}
	dbserver.Db = make(map[string]*sql.DB)
	dbserver.Db["MainDbConnect"] = db

}


func (dbserver *DbServer) AddDbNode(dbName string) {
	psql := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		dbserver.Host, dbserver.Port, dbserver.User, dbserver.Password,dbName)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(90)//avoid tcp or client error
	dbserver.Db[dbName] = db
}






func DbNameExist(db *sql.DB,name string) (bool){
	exist := false
	list := []string{}
	command:= "SELECT datname FROM pg_database WHERE datistemplate = false;"
	result , err := db.Query(command)
	if err != nil {
		log.Println(err)
		return false
	}
	defer result.Close()
	for result.Next() {
		var name string
		if err := result.Scan(&name); err != nil {
			log.Println(err)
			return false
		}
		list = append(list, name)
	}
	if err := result.Err(); err != nil {
		log.Println(err)
		return false

	}
	for _,v := range list{
		if v == name {
			exist = true
		}
	}
	return exist
}

func TableExist(db *sql.DB,name string) (bool){
	exist := false
	list := []string{}
	command:= "SELECT table_name FROM information_schema.tables WHERE table_schema='public'"
	result , err := db.Query(command)
	if err != nil {
		log.Println(err)
		return false

	}
	defer result.Close()
	for result.Next() {
		var name string
		if err := result.Scan(&name); err != nil {
			log.Println(err)
			return false
		}
		list = append(list, name)
	}
	if err := result.Err(); err != nil {
		log.Println(err)
		return false

	}
	for _,v := range list{
		if v == name {
			exist = true
		}
	}
	return exist
}

func CreateUsersTable(db *sql.DB) {
	command := `
    CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    account TEXT UNIQUE NOT NULL,
    password TEXT,
    apiKey TEXT
    );
    `
	_,err := db.Exec(command)
	if err != nil {
		panic(err)
	}
}

func CreateTaskTable(db *sql.DB) {
	command := `
    CREATE TABLE task (
    taskid TEXT PRIMARY KEY,
    time BIGINT,
    apiKey TEXT,
    type TEXT,
    number INT,
    cumulative BOOLEAN,
    parameter TEXT,
    rate NUMERIC
    );
    `
	_,err := db.Exec(command)
	if err != nil {
		panic(err)
	}
}


func CreateStatisticsTable(db *sql.DB) {
	command := `
CREATE TABLE statistics (
    id SERIAL PRIMARY KEY,
	taskid TEXT,
	answer TEXT,
	spendTime INT,
	solverServer TEXT,
	resault BOOLEAN
    );
    `
	_,err := db.Exec(command)
	if err != nil {
		panic(err)
	}
}


//check all db create
func DbAuto(database DbServer)  {
	if !DbNameExist(database.Db["MainDbConnect"],"SolverDb"){
		command := `
    				CREATE DATABASE  "SolverDb"`
		_,err := database.Db["MainDbConnect"].Exec(command)
		if err != nil {
			panic(err)
		}
		db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=74102589 dbname=SolverDb sslmode=disable")
		var times = 0
		for err != nil && times <10{
			db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=74102589 dbname=SolverDb sslmode=disable")
			times += 1
			time.Sleep(10*time.Second)
			log.Println("wait for creat database")
		}
		if err != nil {
			panic(err)
		}
		db.Close()
	}
	database.AddDbNode("SolverDb")
	if !TableExist(database.Db["SolverDb"],"users"){
		CreateUsersTable(database.Db["SolverDb"])
		}
	if !TableExist(database.Db["SolverDb"],"task"){
		CreateTaskTable(database.Db["SolverDb"])
	}
	if !TableExist(database.Db["SolverDb"],"statistics"){
		CreateStatisticsTable(database.Db["SolverDb"])
	}
	}



