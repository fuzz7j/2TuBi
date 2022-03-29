### 2TuBi

t00ls自动签到 ~~+查询域名，每天可获得2TuBi~~，使用Server酱推送。

~~无第三方库，编译更方便。~~

~~使用<http://www.beianw.com/>获取每天最新备案域名。~~

域名查询模块对t00ls造成负担，删除此功能，避免封号。

go version 1.16

### Useage

修改main.go中如下内容，编译后直接运行即可。
```go
username   = ""  //用户名
password   = ""  //密码md5
questionid = "0"  //安全问题id，默认0
answer     = ""  //安全问题答案
sendkey    = "" //Server酱sendkey
```

#### 安全问题id
0 = 没有安全提问  
1 = 母亲的名字  
2 = 爷爷的名字  
3 = 父亲出生的城市  
4 = 您其中一位老师的名字  
5 = 您个人计算机的型号  
6 = 您最喜欢的餐馆名称  
7 = 驾驶执照的最后四位数字