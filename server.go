package main

import (
	"db"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"solver"
	"strconv"
	"strings"
	"time"
)


type AccountType struct {
	Account string `json:"Account"`
	Password string `json:"Password"`
}

type Report struct {
	Answer string `json:"Answer"`
	Correct bool `json:"Correct"`
}




type RedisServer struct {
	Host     string `yaml:"RedisHost"`
	Port     string `yaml:"RedisPort"`
	Password string `yaml:"RedisPassword"`
	Db int `yaml:"RedisDb"`
	Client *redis.Client
}



var (
	redisServer RedisServer
	TaskList  map[string]string
	apikey  solver.ApiKey
	dbserver db.DbServer
)
func (redisServer *RedisServer) Initialize()  {
	redisServer.Client = redis.NewClient(&redis.Options{
		Addr: redisServer.Host+":"+redisServer.Port,
		Password: "",
		DB: redisServer.Db,
	})

}



func main()  {
	//read config
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &apikey)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &dbserver)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &redisServer)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}


	//check db,db init
	dbserver.Initialize()
	db.DbAuto(dbserver)

	redisServer.Initialize()
	go PrepareTask(dbserver.Db["SolverDb"],redisServer.Client,apikey)




	//start server
	r := mux.NewRouter()
	r.HandleFunc("/api/account/new", newAccount).Methods("POST")
	r.HandleFunc("/api/account/{apikey}", accountStatus).Methods("GET")
	r.HandleFunc("/api/task/book", book).Methods("POST")
	r.HandleFunc("/api/task/{taskid}", getTask).Methods("GET")
	r.HandleFunc("/api/report", report).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))

}





//post


//ok
func newAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var acc AccountType
	err := json.NewDecoder(r.Body).Decode(&acc)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"Error":err.Error()})
	}else {
		apikey,err := db.NewAccount(dbserver.Db["SolverDb"],acc.Account,acc.Password)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"Error":err.Error()})
		}else {
			json.NewEncoder(w).Encode(map[string]string{"apikey":apikey})
		}
	}
}


//ok
func book(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request BookingRequest
	var respone BookStatus
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respone.Booking = "Error"
		respone.Status = err.Error()
		json.NewEncoder(w).Encode(respone)
	}else {
		status := request.CheckPostData(dbserver.Db["SolverDb"])
		if status != "" {
			respone.Booking = "Error"
			respone.Status = status
			json.NewEncoder(w).Encode(respone)
		}else {
			status = request.CheckDb(dbserver.Db["SolverDb"])
			if status != "" {
				if status == "Task number update successful" {
					urlParse := strings.Replace(strings.Replace(request.TaskInfo["url"].(string),"https://","https@3A@2F@2F",-1),"/","@2F",-1)
					taskid := request.Key+strconv.FormatInt(request.BookingTime,10)+urlParse
					respone.Booking = "Ok"
					respone.Status = status
					respone.TaskId = taskid
					json.NewEncoder(w).Encode(respone)
				}else {
					respone.Booking = "Error"
					respone.Status = status
					json.NewEncoder(w).Encode(respone)
				}
			}else {
				request.upload(dbserver.Db["SolverDb"])
				urlParse := strings.Replace(strings.Replace(request.TaskInfo["url"].(string),"https://","https@3A@2F@2F",-1),"/","@2F",-1)
				taskid := request.Key+strconv.FormatInt(request.BookingTime,10)+urlParse
				respone.Booking = "Ok"
				respone.Status = "Booking successful"
				respone.TaskId = taskid
				json.NewEncoder(w).Encode(respone)
			}
		}
	}

	/*
		jsonByte, err := json.Marshal(request.TaskInfo)
		if err != nil {
			panic( err)
		}
		fmt.Println(string(jsonByte))
	*/
}


func report(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req Report
	var respone map[string]string
	respone = make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respone["Error"] = err.Error()
	}else {
		command := "UPDATE statistics SET resault="+strconv.FormatBool(req.Correct)+" WHERE answer='"+req.Answer+"';"
		_,err :=dbserver.Db["SolverDb"].Exec(command)
		if err != nil {
			log.Println(err)
			respone["Error"] = err.Error()
		}else {
			respone["Status"] = "ok"
		}
	}
	json.NewEncoder(w).Encode(respone)
}


//get
//ok
func accountStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result := db.SearchTask(dbserver.Db["SolverDb"],params["apikey"])
	json.NewEncoder(w).Encode(&db.Search{Apikey: params["apikey"],List: result})
}

func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	fmt.Println(params["taskid"])
	result := map[string]interface{}{}
	if _,exist := TaskList[params["taskid"]];exist {
		expired := true
		for expired {
			element, err := redisServer.Client.LPop(params["taskid"]).Result()
			fmt.Println(element,err)
			switch err {
			case redis.Nil:
				result = map[string]interface{}{
					"error": "Task is running please retry" ,
				}
				expired = false
			case nil:
				expireTime,err := strconv.ParseInt(element[0:10], 10, 64)
				if err != nil {
					log.Println(err)
					result = map[string]interface{}{
						"Error": err ,
					}
					expired = false
				}
				if time.Now().Unix()+5<expireTime{
					result = map[string]interface{}{
						"Answer": element[10:] ,
					}
					expired = false
				}
			default:
				log.Println(err)
				result = map[string]interface{}{
					"Error": err ,
				}
				expired = false
			}
		}
	}else {
		result = map[string]interface{}{
			"Error": "Task not exist(not start or expired)" ,
		}
	}
	json.NewEncoder(w).Encode(result)
}




















//for test






























