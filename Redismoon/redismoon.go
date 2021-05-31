package Redismoon

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mmcdole/gofeed"
	"strconv"
	"time"
)



var actlastime int64 = 60 * 5
var cachelastime time.Duration = 60 * 3
var FeedLastTime time.Duration = 60 * 20


func Redisinit(){

	//连接redis，创建数据结构
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()
	var ctx = context.Background()

	pong, err := rdb.Ping(ctx).Result()
	fmt.Println(pong, err)

	//gofeed.Feed


	// zset activeusr
	// zadd key score member
	// zadd activeusr timeUnix uid

	timeUnix:=time.Now().Unix()
	initusr := "test5"
	res, err := rdb.Do(ctx, "zadd", "activeusr",timeUnix, initusr).Result()
	//res, err = rdb.Do(ctx, "zadd", "activeusr",timeUnix, 645565).Result()
	fmt.Println(res,err)

	c := new(Cache)
	c.Rssurl = "initurl"
	ok, err := c.SaveInRedis()
	if !ok{
		fmt.Println("init cache err", err)
	}else {
		fmt.Println("init cache success")
	}

}

type UserFeed struct {
	UserName string
	Rssurl string
	Feed gofeed.Feed
}

func (F *UserFeed) SaveRedis() (error) {
	data, err := json.Marshal(F.Feed)
	if err != nil {
		return err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()
	var ctx = context.Background()
	//todo 过期时间
	_ , err = rdb.HSet(ctx, F.UserName, F.Rssurl, data).Result()
	_, err = rdb.Expire(ctx, F.UserName, FeedLastTime*time.Second).Result()
	if err != nil {
		return err
	}
	// fmt.Println("in UserFeed save", res)
	return  nil
}

func (F *UserFeed) GetRedis() (error){
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()
	var ctx = context.Background()

	rebytes,err := rdb.HGet(ctx, F.UserName, F.Rssurl).Result()
	if err != nil {
		return err
	}
	//json反序列化
	object := &gofeed.Feed{}
	err = json.Unmarshal([]byte(rebytes),object)
	if err != nil {
		return err
	}

	F.Feed = *object

	return nil


}


type Cache struct {
	Rssurl string
	Feed gofeed.Feed
}

func (c *Cache) SaveInRedis() (bool, error){
	data, err := json.Marshal(c.Feed)
	if err != nil {
		return false, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()

	var ctx = context.Background()

	res, err := rdb.Set(ctx, c.Rssurl, data, cachelastime * time.Second).Result()
	if err != nil {
		return false, err
	}
	fmt.Println("in saveInRedis", res)
	return true, nil
}


func (c *Cache) GetFromRedis() (bool, error){
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()
	var ctx = context.Background()

	rebytes,_ := rdb.Get(ctx, c.Rssurl).Result()
	//json反序列化
	object := &gofeed.Feed{}
	err := json.Unmarshal([]byte(rebytes),object)
	if err != nil {
		return false, err
	}

	c.Feed = *object

	return true, nil
}


func Getactiveusr() []string{
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()
	var ctx = context.Background()


	timeUnix:=time.Now().Unix()

	//ZRANGEBYSCORE zset min max 显示最小和最大之间的数字
	//ZRANGEBYSCORE activeusr timeUnix-actlastime +inf
	mintime := timeUnix-actlastime

	res, err := rdb.ZRangeByScore(ctx, "activeusr", &redis.ZRangeBy{
		Min:  strconv.FormatInt(mintime, 10),
		Max: "+inf",
		Offset: 0,
		Count: 50,
	}).Result()

	fmt.Println(res, err)
	//删除过期数值
	//ZREMRANGEBYSCORE key min max
	_, err = rdb.Do(ctx, "ZREMRANGEBYSCORE","activeusr","-inf",mintime).Result()


	return res
}

func SetActUser(userName string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()
	var ctx = context.Background()


	timeUnix:=time.Now().Unix()

	_, err := rdb.Do(ctx, "zadd", "activeusr",timeUnix, userName).Result()
	//res, err = rdb.Do(ctx, "zadd", "activeusr",timeUnix, 645565).Result()
	//fmt.Println(res,err)

	return err
}




