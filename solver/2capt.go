package solver

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"request"
	"strconv"
	"strings"
	"time"
)



var (
	ErrTwoApi = errors.New("2captcha: API error")
	ErrTwoTimeout = errors.New("2captcha: Request timeout")
)

//use 2captcha server
type TwoTask struct {
	RequestData map[string]interface{}
	TaskId string
	Result string
	CaptchaTimeout int64
	RequestInterval int
	SpendTime int64
	GenerateTime int64
	Client request.Request
}

func (task *TwoTask) Initialize (key string)  {
	task.RequestData = map[string]interface{}{
		"key":key,
	}
	task.CaptchaTimeout = CaptchaTimeout
	task.RequestInterval = RequestInterval

}

func (task *TwoTask) AddTaskParameters (name string,value interface{}) {
	task.RequestData[name] = value
}

func (task *TwoTask) creatTask () error {
	header := map[string]string{
		"Content-Type":"application/x-www-form-urlencoded",
	}
	data:=url.Values{}
	for k,v := range task.RequestData{
		data.Add(k,v.(string))
	}
	task.SpendTime = time.Now().Unix()
	resp , err := task.Client.Post("http://2captcha.com/in.php",header,data.Encode())
	if err != nil {
		return err
	}
	if strings.HasPrefix(resp.Text , "OK|"){
		task.TaskId = resp.Text[3:]
	}else {
		log.Println(resp.Text)
		return ErrTwoApi

	}
	fmt.Println(resp.Text)
	return nil
}

func (task *TwoTask) waitAnswer() error {
	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	postdata := map[string]string{"key": task.RequestData["key"].(string), "id": task.TaskId,"action":"get"}
	data := url.Values{}
	for k, v := range postdata {
		data.Add(k, v)
	}
	for (time.Now().Unix()-task.SpendTime < task.CaptchaTimeout){
		resp, err := task.Client.Post("http://2captcha.com/res.php", header, data.Encode())
		if err != nil {
			return err
		}
		if resp.Text != "CAPCHA_NOT_READY" {
			if strings.HasPrefix(resp.Text , "OK|"){
				task.Result = resp.Text[3:]
				task.SpendTime = time.Now().Unix()-task.SpendTime
				task.GenerateTime = time.Now().Unix()
				return nil
			}else {
				log.Println(resp.Text)
				return ErrTwoApi
			}
		}
		time.Sleep(time.Duration(task.RequestInterval)*time.Second)
	}
	return ErrTwoTimeout
}

//create task and wait to solve
func (task *TwoTask) Solve(channel chan *Answer) {
	err := task.creatTask()
	if err != nil {
		channel <- &Answer{
			Answer:err.Error(),
			SpendTime:999999,
			SolverServer:"2captcha",
			GenerateTime:0,
		}
	}else {
		err = task.waitAnswer()
		if err != nil {
			channel <- &Answer{
				Answer:err.Error(),
				SpendTime:999999,
				SolverServer:"2captcha",
				GenerateTime:0,
			}
		}else {
			channel <- &Answer{
				Answer:task.Result,
				SpendTime:task.SpendTime,
				SolverServer:"2captcha",
				GenerateTime:task.GenerateTime,
			}
		}
	}
}

//get server balance
func (task *TwoTask) GetBalance() (float64,error) {
	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	postdata := map[string]string{"key": task.RequestData["key"].(string),"action":"getbalance"}
	data := url.Values{}
	for k, v := range postdata {
		data.Add(k, v)
	}
	resp, err := task.Client.Post("http://2captcha.com/res.php", header, data.Encode())
	if err != nil {
		return 0.0,err
	}
	result, err:= strconv.ParseFloat(resp.Text,64)
	if err != nil {
		return 0.0,ErrTwoApi
	}
	return result,nil
}


//for test
func (task *TwoTask) MockTest(channel chan *Answer) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int() % 119
	time.Sleep(time.Duration(n)*time.Second)
	//fmt.Println(time.Now().Unix(),"   ",n,task.RequestData)
	channel <- &Answer{
			Answer:mockAnswer(),
			SpendTime:int64(n),
			SolverServer:"2captcha",
			GenerateTime:time.Now().Unix(),
		}
	}


