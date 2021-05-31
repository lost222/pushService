package main

import (
	"./Redismoon"
	"./model"
	"./myrss"
	"fmt"
	"github.com/mmcdole/gofeed"
	"time"
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
	}

	//todo 修改Feed表中LatesTitle项目

	//for _, rss := range rssUrls{
	//	userWhoSub := model.SearchRecordUser(rss)
	//	//拿到feed
	//	feed, err := myrss.FetchURL(fp, rss)
	//	//直接序列化解决深拷贝问题
	//	if err != nil{
	//		fmt.Println("in FetchURL", err)
	//	}
	//
	//	//推送到这些user的用户订阅流
	//	for _, username := range userWhoSub{
	//		countSetTime++
	//		var uf Redismoon.UserFeed
	//		uf.UserName = username
	//		uf.Rssurl = rss
	//		uf.Feed = *feed
	//		err = uf.SaveRedis()
	//		if err != nil{
	//			fmt.Println("in SaveRedis", err)
	//		}
	//	}
	//}

	fmt.Println("countSetTime=", countSetTime)
}

//var circleTime int = 5 *60


func main() {
	model.InitDB()
	Redismoon.Redisinit()


	startTime := time.Now().UnixNano()

	singleCircle()

	endTime := time.Now().UnixNano()

	seconds:= float64((endTime - startTime) / 1e9)


	fmt.Println("pushService 1 circle run time", seconds, "s")

}
