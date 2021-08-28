package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	action     = "login"
	username   = "" //用户名
	password   = "" //密码md5
	questionid = "0" //安全问题ID，默认0为未设置
	answer     = "" //安全问题答案
	sendkey    = "" //Server酱sendkey
)

type Response struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	Formhash   string `json:"formhash"`
	Mark       string `json:"mark"`
	Cookie     string
	Signsubmit string
}

var r Response
var num int = 0

func main() {
	getCookie, _ := cookiejar.New(nil)
	client := &http.Client{Jar: getCookie}
	resp, _ := client.PostForm("https://www.t00ls.net/login.json", url.Values{"action": {action}, "username": {username}, "password": {password}, "questionid": {questionid}, "answer": {answer}})
	json.NewDecoder(resp.Body).Decode(&r)
	if r.Status != "success" {
		fmt.Println("登陆失败，一小时后重试。")
		time.Sleep(time.Hour)
		main()
	}
	defer resp.Body.Close()
	r.Signsubmit = "true"
	ajaxsign(r, client)
	domainsearch(r, client, getdomain())

}

// t00ls签到
func ajaxsign(r Response, client *http.Client) {
	resp, _ := client.PostForm("https://www.t00ls.net/ajax-sign.json", url.Values{"signsubmit": {r.Signsubmit}, "formhash": {r.Formhash}})
	defer resp.Body.Close()
	var sign Response
	json.NewDecoder(resp.Body).Decode(&sign)
	if sign.Status == "success" {
		fmt.Println("签到成功")
		push(time.Now().Format("2006/01/02 15:04") + "签到成功")
	} else if sign.Message == "alreadysign" {
		fmt.Println("今日已完成签到。")
	} else {
		fmt.Println("签到失败，1小时后重试。")
		time.Sleep(time.Hour)
		ajaxsign(r, client)
	}
}

// 方糖推送
func push(msg string) {
	url := "https://sctapi.ftqq.com/" + sendkey + ".send?title=t00ls&desp=" + url.QueryEscape(msg)
	http.Get(url)
}

// 获取最新备案域名。
func getdomain() []string {
	var res []string
	resp, err := http.Get("http://www.beianw.com/")
	if err != nil {
		fmt.Println("获取失败")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	varhrefRegexp := regexp.MustCompile("\\w{0,62}\\.com")
	match := varhrefRegexp.FindAllString(string(body), -1)
	for i := 0; i < len(match); i++ {
		flag := false
		for j := i + 1; j < len(match); j++ {
			if match[i] == match[j] {
				flag = true
				break
			}
		}
		if !flag {
			res = append(res, match[i])
		}
	}
	return res
}

// 查询域名并查询tubi获取日志，如果包含域名则查询成功。
func domainsearch(r Response, client *http.Client, res []string) {
	for i := 1; i < len(res); i++ {
		client.PostForm("https://www.t00ls.net/domain.html", url.Values{"domain": {res[i]}, "formhash": {r.Formhash}, "querydomainsubmit": {"%E6%9F%A5%E8%AF%A2"}})
		tubilog, err := client.Get("https://www.t00ls.net/members-tubilog.json")
		if err != nil {
			//
		}
		defer tubilog.Body.Close()
		body, err := io.ReadAll(tubilog.Body)
		if strings.Contains(string(body), res[i]) == true {
			fmt.Printf("%s 域名查询成功，Tubi Get！", res[i])
			push(time.Now().Format("2006/01/02 15:04") + res[i] + "域名查询成功")
			num = 0
			break
		} else {
			num++
		}
	}
	if num != 0 && num < 100 {
		time.Sleep(time.Hour * 2)
		domainsearch(r, client, getdomain())
	}
}
