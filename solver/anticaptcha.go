package solver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"request"
	"strconv"
	"strings"
	"time"
)


var (
	ErrAntiApi = errors.New("anticaptcha: API error")
	ErrAntiTimeout = errors.New("anticaptcha: Request timeout")
)

//use anticaptcha server
type AntiTask struct {
	RequestData map[string]interface{}
	TaskId string
	Result string
	CaptchaTimeout int64
	RequestInterval int
	SpendTime int64
	GenerateTime int64
	Client request.Request

}

func (task *AntiTask) Initialize (key string)  {
	task.RequestData = map[string]interface{}{
		"clientKey":key,
		"task": map[string]interface{}{
		},
	}
	task.CaptchaTimeout = CaptchaTimeout
	task.RequestInterval = RequestInterval
}

func (task *AntiTask) AddTaskParameters (name string,value interface{}) {
	task.RequestData["task"].(map[string]interface{})[name] = value
}

func (task *AntiTask) creatTask () error {
	header := map[string]string{
		"Content-Type": "application/json",
	}
	jsonByte, err := json.Marshal(task.RequestData)
	if err != nil {
		return err
	}
	jsonString := request.Byte2String(jsonByte)
	task.SpendTime = time.Now().Unix()
	resp , err := task.Client.Post("https://api.anti-captcha.com/createTask",header,jsonString)
	if err != nil {
		return err
	}
	if strings.HasPrefix(resp.Text , "{\"errorId\":0,"){
		task.TaskId = resp.Text[22:len(resp.Text)-1]
	}else {
		log.Println(resp.Text)
		return ErrAntiApi
		}
	fmt.Println(resp.Text)
	return nil
}

//todo can add other type captcha in switch
func (task *AntiTask) waitAnswer() error {
	header := map[string]string{
		"Content-Type": "application/json",
	}
	data := map[string]string{"clientKey": task.RequestData["clientKey"].(string), "taskId": task.TaskId}
	jsonByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	jsonString := request.Byte2String(jsonByte)

	for (time.Now().Unix()-task.SpendTime < task.CaptchaTimeout) {
		resp, err := task.Client.Post("https://api.anti-captcha.com/getTaskResult", header, jsonString)
		if err != nil {
			return err
		}
		fmt.Println(resp.Text)
		var ResponseJson interface{}
		err = json.Unmarshal([]byte(resp.Text), &ResponseJson)
		if err != nil {
			return err
		}
		if ResponseJson.(map[string]interface{})["errorId"] == 0.0 {
			if ResponseJson.(map[string]interface{})["status"] == "ready" {
				switch task.RequestData["task"].(map[string]interface{})["type"].(string) {
				case "ImageToTextTask":
					task.Result = ResponseJson.(map[string]interface{})["solution"].(map[string]interface{})["text"].(string)
				default:
					task.Result = ResponseJson.(map[string]interface{})["solution"].(map[string]interface{})["gRecaptchaResponse"].(string)
				}
				task.SpendTime = time.Now().Unix()-task.SpendTime
				task.GenerateTime = time.Now().Unix()
				return nil
			} else {
				time.Sleep(time.Duration(task.RequestInterval) * time.Second)
			}
		} else {
			log.Println(resp.Text)
			return ErrAntiApi
		}
	}
	return ErrAntiTimeout
}

//create task and wait to solve
func (task *AntiTask) Solve(channel chan *Answer) {
	err := task.creatTask()
	if err != nil {
		channel <- &Answer{
			Answer:err.Error(),
			SpendTime:999999,
			SolverServer:"anticaptcha",
			GenerateTime:0,
		}
	}else {
		err = task.waitAnswer()
		if err != nil {
			channel <- &Answer{
				Answer:err.Error(),
				SpendTime:999999,
				SolverServer:"anticaptcha",
				GenerateTime:0,
			}
		}else {
			channel <- &Answer{
				Answer:task.Result,
				SpendTime:task.SpendTime,
				SolverServer:"anticaptcha",
				GenerateTime:task.GenerateTime,
			}
		}
	}
}

//get server balance
func (task *AntiTask) GetBalance() (float64,error) {
	header := map[string]string{
		"Content-Type": "application/json",
	}
	data := map[string]string{"clientKey": task.RequestData["clientKey"].(string)}
	jsonByte, err := json.Marshal(data)
	if err != nil {
		return 0.0,err
	}
	jsonString := request.Byte2String(jsonByte)
	resp, err := task.Client.Post("https://api.anti-captcha.com/getBalance", header, jsonString)
	if err != nil {
		return 0.0,err
	}
	if strings.HasPrefix(resp.Text , "{\"errorId\":0,") {
		result, err := strconv.ParseFloat(resp.Text[23:len(resp.Text)-1], 64)
		if err != nil {
			return 0.0, err
		}
		return result,nil
	}else {
		return 0.0,ErrAntiApi
	}
}

//for test
func (task *AntiTask) MockTest(channel chan *Answer) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int() % 120
	time.Sleep(time.Duration(n)*time.Second)
	//fmt.Println(time.Now().Unix(),"   ",n,task.RequestData)
	channel <- &Answer{
		Answer:mockAnswer(),
		SpendTime:int64(n),
		SolverServer:"anticaptcha",
		GenerateTime:time.Now().Unix(),
	}
}
