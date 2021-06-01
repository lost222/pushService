package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
)

const (
	consulAddress = "127.0.0.1:8500"
	localIp       = "127.0.0.1"
	localPort     = 81
)

const PushServiceID = "666543"


func ConsulRegister()  {
	// 创建连接consul服务配置
	config := consulapi.DefaultConfig()
	config.Address = consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("consul client error : ", err)
	}

	// 创建注册到consul的服务到
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = PushServiceID
	registration.Name = "PushService"
	registration.Port = localPort
	registration.Tags = []string{"PushService"}
	registration.Address = localIp

	// 增加consul健康检查回调函数
	check := new(consulapi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d", registration.Address, registration.Port)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "30s" // 故障检查失败30s后 consul自动将注册服务删除
	registration.Check = check

	// 注册服务到consul
	err = client.Agent().ServiceRegister(registration)
}


