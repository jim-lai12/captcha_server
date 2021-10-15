package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
)

type BookingRequest struct {
	BookingTime int64                 //use unix time
	Key string                        //user Key
	Type string                       //task type "recaptchav2","recaptchav3"
	Number int                        //the number of captcha you want to book
	Cumulative bool                   //if true you can increase the booking captcha number by post many time.
	TaskInfo map[string]interface{}        //task parameter
	Rate float64                      //not necessary(has default value),the rate of captcha number.when captcha is error,you can change new answer fast
}

type BookStatus struct {
	Booking string
	Status string
	TaskId string
}



//check user post data
func (request BookingRequest) CheckPostData(db *sql.DB) string {
	if time.Now().Unix() > request.BookingTime {
		return "Booking time Error.It is past time"
	}else if time.Now().Unix()+86400 < request.BookingTime {
		return "Booking time Error.It is bigger than 24h"
	}
	if request.Key == ""{
		return "Key Error."
	}else{
		//var aaa string
		command:= "select * from users where apikey='"+ request.Key +"'"
		result , _ := db.Query(command)
		defer result.Close()
		if !result.Next(){
			return "Key Error.No this Key"
		}
	}
	if request.Type == "recaptchav2"{
		toCheckList := []string {"url","websiteKey"}
		resault := []int{}
		for i,value := range toCheckList{
			if !inTaskInfo(request.TaskInfo,value) {
				resault = append(resault, i)
			}
		}
		if len(resault)!=0{
			errormesssage := "TaskInfo Error.No "
			for _,value:= range resault{
				errormesssage += toCheckList[value] + " "
			}
			errormesssage += "parameter"
			return errormesssage
		}
	}else if request.Type == "recaptchav3" {
		toCheckList := []string {"url","websiteKey","minScore"}
		resault := []int{}
		for i,value := range toCheckList{
			if !inTaskInfo(request.TaskInfo,value) {
				resault = append(resault, i)
			}
		}
		if len(resault)!=0{
			errormesssage := "TaskInfo Error.No "
			for _,value:= range resault{
				errormesssage += toCheckList[value] + " "
			}
			errormesssage += "parameter"
			return errormesssage
		}
	}else {
		return "Type Error.Not support this type"
	}

	return ""
}

/*
   dbname : task
   taskid :Key+BookingTime+TaskInfo.url(replace https://,/)
   time : BookingTime
   apiKey : Key
   type : Type
   number : Number
   cumulative : Cumulative
   parameter : parameter json string
   rate : Rate
*/

//check db data whether repeat
//todo check row type and is sentment
func (request BookingRequest) CheckDb(db *sql.DB) string {
	urlParse := strings.Replace(strings.Replace(request.TaskInfo["url"].(string),"https://","https@3A@2F@2F",-1),"/","@2F",-1)
	taskid := request.Key+strconv.FormatInt(request.BookingTime,10)+urlParse
	command:= "select cumulative from task where taskid='"+ taskid +"'"
	result , _ := db.Query(command)
	defer result.Close()
	if result.Next(){
		var row bool
		if err := result.Scan(&row); err != nil {
			log.Fatal(err)
		}
		if  row {
			command = "UPDATE task SET number= number + "+strconv.Itoa(request.Number)+" WHERE taskid='"+taskid+"'"
			db.Exec(command)
			return "Task number update successful"
		}
		return "Task existed and can't Cumulative"
	}
	return ""
}


//upload to Db
func (request BookingRequest) upload(db *sql.DB) string {
	urlParse := strings.Replace(strings.Replace(request.TaskInfo["url"].(string),"https://","https@3A@2F@2F",-1),"/","@2F",-1)
	taskid := request.Key+strconv.FormatInt(request.BookingTime,10)+urlParse
	jsonbyte , err := json.Marshal(request.TaskInfo)
	jsonString := string(jsonbyte)

	command1 := `
	INSERT INTO task (taskid,time,apiKey,type,number,cumulative,parameter,rate)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`
	_,err = db.Exec(command1,taskid,request.BookingTime,request.Key,request.Type,request.Number,request.Cumulative,jsonString,request.Rate)
	if err != nil {
		return err.Error()
	}
	return ""
}

func inTaskInfo(mapObject map[string]interface{},value string) (bool) {
	if _,exist := mapObject[value]; exist{
		return true
	}
	return false
}
