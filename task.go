package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"math"
	"solver"
	"strconv"
	"time"
)

const defultRate = 1.1

type Task struct {
	Context context.Context
	Cancel context.CancelFunc
	Function func(ctx context.Context,db *sql.DB,redis *redis.Client,taskid string,bookTime int64,captchaType string,number int,parameter string,rate float64,key solver.ApiKey)
	Booktime int64
}

func PrepareTask(db *sql.DB,redisServer *redis.Client,key solver.ApiKey)  {
	taskMap := make(map[string]Task)
	//taskMap := make(map[string]func(db *sql.DB,redis *redis.Client,taskid string,bookTime int64,captchaType string,number int,parameter string,rate float64,key solver.ApiKey))
	//taskContextMap := make(map[string]context.Context)
	//mainContext := context.Background()
	var taskid,captchaType,parameter string
	var bookTime int64
	var number int
	var rate float64
	for true{
		command:= "select taskid,time,type,number,parameter,rate from task where time > " +strconv.FormatInt(time.Now().Unix(),10)+"and time <"+strconv.FormatInt(time.Now().Unix()+130,10)
		result , _ := db.Query(command)
		defer result.Close()
		for result.Next(){
			if err := result.Scan(&taskid,&bookTime,&captchaType,&number,&parameter,&rate); err != nil {
				log.Println(err)
			}
			if _,exist := taskMap[taskid]; !exist{
				fmt.Println(taskid,"not exist")
				//ctx,cancel := context.WithCancel(mainContext)
				ctx,cancel := context.WithCancel(context.Background())
				taskMap[taskid] = Task{Context:ctx,Cancel:cancel,Function: SubTask,Booktime: bookTime}
				TaskList[taskid] = "start"
				go taskMap[taskid].Function(taskMap[taskid].Context,db,redisServer,taskid,bookTime,captchaType,number,parameter,rate,key)
				//taskContextMap[taskid],_ = context.WithCancel(mainContext)
				//taskMap[taskid] = SubTask
				//go taskMap[taskid](db,redisServer,taskid,bookTime,captchaType,number,parameter,rate,key)
			}
		}
		for k,v := range taskMap{
			if time.Now().Unix() >= v.Booktime + 180 {
				fmt.Println("Stop Task : ",k,time.Now().Unix(),v.Booktime)
				taskMap[k].Cancel()
				delete(taskMap,k)
				delete(TaskList,k)
				redisServer.Del(k)
			}
		}
		time.Sleep(10*time.Second)
	}
}

func SubTaskMock(ctx context.Context,db *sql.DB,redis *redis.Client,taskid string,bookTime int64,captchaType string,number int,parameter string,rate float64,key solver.ApiKey)  {
	/*
		fmt.Println(taskid)
		fmt.Println(bookTime)
		fmt.Println(number)
		fmt.Println(parameter)
		fmt.Println(rate)
	*/
	if rate != 0{
		number = int(math.Floor(float64(number)*rate))
	}else {
		number = int(math.Floor(float64(number)*defultRate))
	}
	fmt.Println(number)
	solver.StartSolverMock(ctx,db,redis,taskid,number,captchaType,parameter,&key)
}
func SubTask(ctx context.Context,db *sql.DB,redis *redis.Client,taskid string,bookTime int64,captchaType string,number int,parameter string,rate float64,key solver.ApiKey)  {
	/*
		fmt.Println(taskid)
		fmt.Println(bookTime)
		fmt.Println(number)
		fmt.Println(parameter)
		fmt.Println(rate)
	*/
	if rate != 0{
		number = int(math.Floor(float64(number)*rate))
	}else {
		number = int(math.Floor(float64(number)*defultRate))
	}
	fmt.Println(number)
	solver.StartSolver(ctx,db,redis,taskid,number,captchaType,parameter,&key)
}
