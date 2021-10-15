package db

import (
	"database/sql"
	"log"
	"strconv"
	"time"
)

func PrepareTask(db *sql.DB)  {
	for true{
		command:= "select taskid from task where time > " +strconv.FormatInt(time.Now().Unix(),10)+"and time <"+strconv.FormatInt(time.Now().Unix()+120,10)
		result , _ := db.Query(command)
		defer result.Close()
		for result.Next(){
			var taskid string
			if err := result.Scan(&taskid); err != nil {
				log.Println(err)
			}
			println(taskid)
		}
		time.Sleep(10*time.Second)
	}
}
