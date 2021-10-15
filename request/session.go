package request

import (
	"cookiejar"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)



//keep session cookie will auto update
type Session struct {
	Client http.Client
	Jar *cookiejar.Jar
}
//init session
func (session *Session)New() {
	session.Jar,_ = cookiejar.New(nil)
	session.Client.Jar = session.Jar

	tr := &http.Transport{
		MaxIdleConns: 100,
		IdleConnTimeout: 30 * time.Second,
		DisableCompression: false,
	}
	session.Client.Transport = tr
	session.Client.Timeout = 10 * time.Second
}

func (session *Session)SetProxy(ip string,port string,username string,password string)(error) {
	proxy,_:= url.Parse("http://96.2.228.18:8080")//http://96.2.228.18:8080//socks5://104.238.66.161:31337
	session.Client.Transport.(*http.Transport).Proxy = http.ProxyURL(proxy)
	return nil
}


//can load and save MozillaCookieJar file from python or for python
func (session *Session)LoadCookie(FileName string)  {
	ByteText, _ := ioutil.ReadFile(FileName)
	re := regexp.MustCompile(`.*\n`)
	out:=re.FindAll([]byte(ByteText), -1)
	for  i:=4; i<len(out);i++  {
		cookieval:=strings.Split(Byte2String(out[i]),"	")
		var secure bool
		secure,_ = strconv.ParseBool(cookieval[3])
		expires,_ := strconv.ParseInt(cookieval[4], 10, 64)
		cookie := &http.Cookie{
			Name : cookieval[5],
			Value : strings.Replace(cookieval[6],"\n","",-1),
			Path : cookieval[2],
			Domain : cookieval[0],
			Expires : time.Unix(expires, 0),
			Secure : secure,
		}
		u,_:= url.Parse("https://"+cookieval[0])//url.URL{Host: ""}
		session.Client.Jar.SetCookies(u, []*http.Cookie{cookie})
	}
}

func (session *Session)SaveCookie(FileName string)  {
	data := "# Netscape HTTP Cookie File\n# http://curl.haxx.se/rfc/cookie_spec.html\n# This is a generated file!  Do not edit.\n\n"
	for key := range session.Jar.Entries{
		for _,value := range session.Jar.Entries[key]{
			data += "." + value.Domain + "\tTRUE\t" + value.Path + "\t" + strings.ToUpper(strconv.FormatBool(value.Secure)) + "\t" + strconv.FormatInt(value.Expires.Unix(),10) + "\t" + value.Name + "\t" + value.Value + "\n"
		}
	}
	ioutil.WriteFile(FileName,[]byte(data),0777)
}

func (session *Session)Get(url string ,header map[string]string)(*Response, error){
	req, err := http.NewRequest("GET", url, nil)
	if err !=nil{
		return nil,err
	}
	for k,v := range header{
		req.Header.Add(k,v)
	}
	resp, err := session.Client.Do(req)
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

func (session *Session)Post(url string ,header map[string]string,data string)(*Response, error){
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err !=nil{
		return nil,err
	}
	for k,v := range header{
		req.Header.Add(k,v)
	}
	resp, err := session.Client.Do(req)
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





