package solver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	RequestInterval int =  5
	CaptchaTimeout int64 = 120
)
type ApiKey struct {
	AntiCaptcha []string `yaml:"anti-captcha"`
	TwoCaptcha []string `yaml:"2captcha"`
}


type Solver interface {
	Initialize(key string)
	AddTaskParameters(name string,value interface{})
	Solve(channel chan *Answer)
	GetBalance() (float64,error)
	MockTest(channel chan *Answer)
}



//for mock test
func StartSolverMock(ctx context.Context,psql *sql.DB,redis *redis.Client,taskid string,number int,captchaType string,payload string,key *ApiKey)  {
	var queue AnswerQueue
	queue.Initialize(number)
	rand.Seed(time.Now().Unix())
	solverlist := []Solver{}
	var antiData map[string]interface{}
	var twoData map[string]interface{}
	if captchaType == "recaptchav2" {
		var recaptcha RecaptchaV2
		err := json.Unmarshal([]byte(payload), &recaptcha)
		if err != nil {
			log.Println(err)
		}
		antiData = recaptcha.Parse2captcha()
		twoData = recaptcha.ParseAntiCaptcha()
	}else if captchaType == "recaptchav3" {
		var recaptcha RecaptchaV3
		err := json.Unmarshal([]byte(payload), &recaptcha)
		if err != nil {
			log.Println(err)
		}
		antiData = recaptcha.Parse2captcha()
		twoData = recaptcha.ParseAntiCaptcha()
	}
	for i:=1;i<=number;i+=1{
		r := i%2
		switch r {
		case 1:
			worker := new(AntiTask)
			n := rand.Int() % len(key.AntiCaptcha)
			worker.Initialize(key.AntiCaptcha[n])
			for k,v := range antiData{
				worker.AddTaskParameters(k,v)
			}
			solverlist = append(solverlist, worker)
		case 0:
			worker := new(TwoTask)
			n := rand.Int() % len(key.TwoCaptcha)
			worker.Initialize(key.TwoCaptcha[n])
			for k,v := range twoData{
				worker.AddTaskParameters(k,v)
			}
			solverlist = append(solverlist, worker)
		}
	}
	channel := make(chan *Answer,number)
	fmt.Println(taskid,channel)
	for _,v := range solverlist{
		go v.MockTest(channel)
	}
	go queue.AllPushAndRedis(channel,redis,taskid)
	go queue.CheckItemPop(ctx)
	go queue.SyncDb(ctx,psql,taskid)
	for {
		select {
		case <-ctx.Done():
			close(channel)
			fmt.Println("stop",channel,time.Now().Unix(),taskid)
			return
		default:
			time.Sleep(2*time.Second)
		}
	}
}


//real function
func StartSolver(ctx context.Context,psql *sql.DB,redis *redis.Client,taskid string,number int,captchaType string,payload string,key *ApiKey)  {
	var queue AnswerQueue
	queue.Initialize(number)
	rand.Seed(time.Now().Unix())
	solverlist := []Solver{}
	var antiData map[string]interface{}
	var twoData map[string]interface{}
	if captchaType == "recaptchav2" {
		var recaptcha RecaptchaV2
		err := json.Unmarshal([]byte(payload), &recaptcha)
		if err != nil {
			log.Println(err)
		}
		antiData = recaptcha.Parse2captcha()
		twoData = recaptcha.ParseAntiCaptcha()
	}else if captchaType == "recaptchav3" {
		var recaptcha RecaptchaV3
		err := json.Unmarshal([]byte(payload), &recaptcha)
		if err != nil {
			log.Println(err)
		}
		antiData = recaptcha.Parse2captcha()
		twoData = recaptcha.ParseAntiCaptcha()
	}
	for i:=1;i<=number;i+=1{
		r := i%2
		switch r {
		case 1:
			worker := new(AntiTask)
			n := rand.Int() % len(key.AntiCaptcha)
			worker.Initialize(key.AntiCaptcha[n])
			for k,v := range antiData{
				worker.AddTaskParameters(k,v)
			}
			solverlist = append(solverlist, worker)
		case 0:
			worker := new(TwoTask)
			n := rand.Int() % len(key.TwoCaptcha)
			worker.Initialize(key.TwoCaptcha[n])
			for k,v := range twoData{
				worker.AddTaskParameters(k,v)
			}
			solverlist = append(solverlist, worker)
		}
	}
	channel := make(chan *Answer,number)
	fmt.Println(taskid,channel)
	for _,v := range solverlist{
		go v.Solve(channel)
	}
	go queue.AllPushAndRedis(channel,redis,taskid)
	go queue.CheckItemPop(ctx)
	go queue.SyncDb(ctx,psql,taskid)
	for {
		select {
		case <-ctx.Done():
			close(channel)
			fmt.Println("stop",channel,time.Now().Unix(),taskid)
			return
		default:
			time.Sleep(2*time.Second)
		}
	}
}














type AnswerQueue struct {
	mux sync.RWMutex
	TaskId string
	top int
	end int
	Answers []*Answer
}

type Answer struct {
	Answer string
	SpendTime int64
	SolverServer string
	GenerateTime int64
}

func (queue *AnswerQueue) Initialize(TaskNumber int) {
	queue.Answers = make([]*Answer,TaskNumber)
}

func (queue *AnswerQueue) Push(answer *Answer) {
	queue.mux.Lock()
	queue.Answers[queue.end] = answer
	queue.end +=1
	queue.mux.Unlock()
}

func (queue *AnswerQueue) Pop() (Answer){
	if queue.Isempty(){
		return Answer{}
	}
	queue.mux.Lock()
	defer queue.mux.Unlock()
	queue.top += 1
	return *queue.Answers[queue.top-1]
}

//get result and upload to redis
func (queue *AnswerQueue) AllPushAndRedis(c chan *Answer,redis *redis.Client,taskid string) {
	var now int
	var value string
	var expire int64
	for answer := range c {
		queue.mux.Lock()
		queue.Answers[queue.end] = answer
		now = queue.end
		value = queue.Answers[now].Answer
		expire = queue.Answers[now].GenerateTime + 120
		queue.end += 1
		queue.mux.Unlock()

		go func() {
			err := redis.RPush(taskid,strconv.FormatInt(expire,10)+value).Err()
			//err := redis.Set(taskid+strconv.Itoa(now), value , 120*time.Second).Err()
			if err != nil {
				log.Println(err)
			}
		}()
	}
	fmt.Println("AllPushAndRedis not running")
}

//check channel status and pop item
func (queue *AnswerQueue) CheckItemPop(ctx context.Context) {
	for{
		select {
		case <-ctx.Done():
			fmt.Println("CheckItemPop not running")
			return
		default:
			for !queue.Isempty() {
				queue.mux.Lock()
				if queue.Answers[queue.top].GenerateTime+120 <= time.Now().Unix() {
					queue.top += 1
					queue.mux.Unlock()
				} else {
					queue.mux.Unlock()
				}
			}
		}
	}
}

//check queue status
func (queue *AnswerQueue) Isempty()(bool) {
	queue.mux.Lock()
	defer queue.mux.Unlock()
	if queue.end == queue.top {
		return true
	}
	return false
}

//upload resault to db
func (queue *AnswerQueue) SyncDb(ctx context.Context,psql *sql.DB,taskid string){
	lastSync := 0
	for{
		select {
		case <-ctx.Done():
			fmt.Println("SyncDb not running")
			return
		default:
			command := `INSERT INTO statistics (taskid, answer, spendTime, solverServer)
				VALUES `
			if queue.top != 0 {
				queue.mux.RLock()
				if lastSync < queue.top {
					for lastSync < queue.top {
						command += "('" + taskid + "','" + queue.Answers[lastSync].Answer + "','" + strconv.FormatInt(queue.Answers[lastSync].SpendTime, 10) + "','" + queue.Answers[lastSync].SolverServer + "'),"
						lastSync += 1
					}
					command = command[0:len(command)-1] + ";"
					_, err := psql.Exec(command)
					if err != nil {
						log.Println(err)
					}
					queue.mux.RUnlock()
				}else {
					queue.mux.RUnlock()
				}

			}
			time.Sleep(12*time.Second)
			}
	}
}










//for mock answer
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-"

func mockAnswer() string {
	rand.Seed(time.Now().UnixNano())
	ans := ""
	for i:=0;i<67;i++{
		n := rand.Int() % len(charset)
		ans += string(charset[n])
	}
	return ans
}

