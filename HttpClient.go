package main

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"net/http"
)


const (
	consulAgentAddress = "127.0.0.1:8500"
)

// 从consul中发现服务
func ConsulFindServerHTTP(ServiceName string)  (string, int){
	// 创建连接consul服务配置
	config := consulapi.DefaultConfig()
	config.Address = consulAgentAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("consul client error : ", err)
	}

	// 获取指定service
	//service, _, err := client.Agent().Service("337", nil)
	//if err == nil{
	//	fmt.Println(service.Address)
	//	fmt.Println(service.Port)
	//}

	//只获取健康的service
	serviceHealthy, _, err := client.Health().Service(ServiceName, "", true, nil)
	if err == nil{
		fmt.Println(serviceHealthy[0].Service.Address)
	}

	//多个要考虑负载均衡
	ser := serviceHealthy[0]
	ip := ser.Service.Address
	port := ser.Service.Port
	return ip, port
}

func main()  {
	ip, port := ConsulFindServerHTTP("PushService")
	addr := fmt.Sprintf("http://%v:%v/push", ip, port)
	fmt.Println(addr)

	res,_ :=http.Get(addr)
	defer res.Body.Close()
	fmt.Print(res.Body)
}