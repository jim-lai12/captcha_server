package request

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"time"
	"unsafe"
)

type Request struct {
	Client http.Client
}

//Text is respone bosy string
type Response struct {
	Text string
	Detail *http.Response
}

func Byte2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//set timeout for client object
func SetTimeOut(client *http.Client,timeoutSec int)  {
	client.Timeout=time.Duration(timeoutSec)*time.Second
}


func (request *Request)SetProxy(ip string,port string,username string,password string)(error) {
	proxy,_:= url.Parse("http://96.2.228.18:8080")//http://96.2.228.18:8080//socks5://104.238.66.161:31337
	request.Client.Transport.(*http.Transport).Proxy = http.ProxyURL(proxy)
	return nil
}

func (request *Request)Get(url string ,header map[string]string)(*Response, error){
	req, err := http.NewRequest("GET", url, nil)
	if err !=nil{
		return nil,err
	}
	for k,v := range header{
		req.Header.Add(k,v)
	}
	resp, err := request.Client.Do(req)
	defer resp.Body.Close()
	var body []byte
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		bodyraw, err := gzip.NewReader(resp.Body)
		if err !=nil{
			return nil,err
		}
		body, err = ioutil.ReadAll(bodyraw)
		if err !=nil{
			return nil,err
		}
	default:
		body, err = ioutil.ReadAll(resp.Body)
		if err !=nil{
			return nil,err
		}
	}
	return &Response{Text: Byte2String(body),Detail: resp},nil
}

func (request *Request)Post(url string ,header map[string]string,data string)(*Response, error){
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err !=nil{
		return nil,err
	}
	for k,v := range header{
		req.Header.Add(k,v)
	}
	resp, err := request.Client.Do(req)
	if err !=nil{
		return nil,err
	}
	defer resp.Body.Close()
	var body []byte

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		bodyraw, err := gzip.NewReader(resp.Body)
		if err !=nil{
			return nil,err
		}
		body, err = ioutil.ReadAll(bodyraw)
		if err !=nil{
			return nil,err
		}
	default:

		body, err = ioutil.ReadAll(resp.Body)
		if err !=nil{
			return nil,err
		}

	}
	return &Response{Text: Byte2String(body),Detail: resp},nil
}

/*
consider using to read resp.Body

buffer := bytes.NewBuffer(make([]byte, 4096))
_, err := io.Copy(buffer, request.Body)
if err !=nil{
    return nil, err
}


*/





































func Get(url string ,client *http.Client,header []string)(*http.Response, error){
	start := time.Now()
	req, err := http.NewRequest("GET", url, nil)
	if err !=nil{
		return nil,err
	}
	for  i:=0; i<len(header)/2;i++ {
		req.Header.Add(header[2*i],header[2*i+1])
	}
	resp, err := client.Do(req)
	fmt.Println(time.Until(start))
	return resp,nil
}




