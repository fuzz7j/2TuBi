package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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
	resp, _ := client.PostForm("https://www.t00ls.cc/login.json", url.Values{"action": {action}, "username": {username}, "password": {password}, "questionid": {questionid}, "answer": {answer}})
	json.NewDecoder(resp.Body).Decode(&r)
	if r.Status != "success" {
		fmt.Println("登陆失败，一小时后重试。")
		time.Sleep(time.Hour)
		main()
	}
	defer resp.Body.Close()
	r.Signsubmit = "true"
	ajaxsign(r, client)
}

// t00ls签到
func ajaxsign(r Response, client *http.Client) {
	resp, _ := client.PostForm("https://www.t00ls.cc/ajax-sign.json", url.Values{"signsubmit": {r.Signsubmit}, "formhash": {r.Formhash}})
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