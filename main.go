package main

import (
	"./Redismoon"
	"./consul"
	"./model"
	"./myrss"
	"fmt"
	"github.com/mmcdole/gofeed"
	"net/http"
)



func distinct(s []string) []string{
	var result []string
	set := make(map[string]struct{})
	for _ , item := range s{
		if _,ok := set[item]; !ok{
			set[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}


func singleCircle(){
	fmt.Println("in singleCircle")
	activeUsers := Redismoon.Getactiveusr()
	fmt.Println("activeUsers:" , activeUsers)

	//这个应该是个set
	var rssUrls []string
	for _, user := range activeUsers{
		rss := model.SearchUserSubRecord(user)
		rssUrls = append(rssUrls, rss...)
	}
	rssUrls = distinct(rssUrls)

	fmt.Println("len of rssUrls", len(rssUrls))

	fp := gofeed.NewParser()

	var countSetTime int = 0
	for _, rss := range rssUrls{
		feed, err := myrss.FetchURL(fp, rss)
		if err != nil{
			fmt.Println(err)
			continue
		}
		cc := Redismoon.Cache{
			Rssurl: rss,
			Feed: *feed,
		}
		cc.SaveInRedis()
		countSetTime++
	}

	//todo 修改Feed表中LatesTitle项目


	fmt.Println("countSetTime=", countSetTime)
}

//var circleTime int = 5 *60


func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello From PushService"))
}


func PushHandler(w http.ResponseWriter, r *http.Request) {
	singleCircle()
	w.Write([]byte("PUSH SUCCESS"))
}

func main() {
	model.InitDB()
	Redismoon.Redisinit()
	consul.ConsulRegister()

	http.HandleFunc("/", Handler)
	http.HandleFunc("/push", PushHandler)

	err := http.ListenAndServe("0.0.0.0:81", nil)
	if err != nil {
		fmt.Println("error: ", err.Error())
	}


}
