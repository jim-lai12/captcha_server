package db

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"strings"
	"time"
)

type Search struct {
	Apikey string
	List []*SearchResult
}
type SearchResult struct {
	ApiKey string
	BookingTime string
	Url string
	Number int
}

func NewAccount(db *sql.DB,account string,password string) (string ,error) {
	command1 := `
	INSERT INTO users (account, password, apiKey)
	VALUES ($1, $2 ,$3)
	`
	sha := sha256.New()
	sha.Write([]byte(password))
	passwordHash := sha.Sum(nil)
	apiKey,_ := GenerateRandomString(32)
	_,err := db.Exec(command1,account,hex.EncodeToString(passwordHash),apiKey)
	if err != nil {
		return "", err
	}
	return apiKey,nil
}


func SearchTask(db *sql.DB,apiKey string) []*SearchResult{
	var resultlist = []*SearchResult{}
	command:= "select taskid,number from task where apiKey='"+ apiKey +"' and time > " +strconv.FormatInt(time.Now().Unix(),10)
	result , _ := db.Query(command)
	defer result.Close()
	for result.Next(){
		var taskid string
		var n int
		if err := result.Scan(&taskid,&n); err != nil {
			log.Println(err)
		}
		url := strings.Replace(strings.Replace(taskid[42:],"https%3A%2F%2F","https://",-1),"%2F","/",-1)

		resultlist = append(resultlist, &SearchResult{ApiKey: taskid[:32],BookingTime: taskid[32:42],Url: url,Number: n})
	}

return resultlist
}