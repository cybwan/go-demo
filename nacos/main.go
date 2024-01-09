package main

// https://mdnice.com/writing/6c875ef6c97a4ccc9dbd64a735b59faa

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	// 创建一个Nacos客户端
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "localhost",
			Port:   8848,
		},
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}

	//client.DeregisterInstance(vo.DeregisterInstanceParam{
	//	Ip:          "10.0.0.14",
	//	Port:        8848,
	//	ServiceName: "demo.go",
	//	GroupName:   "group-b",
	//})
	//
	success, err := client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "10.0.0.14",
		Port:        8848,
		ServiceName: "bookwarehouse",
		Weight:      10,
		ClusterName: "cluster-b",
		GroupName:   "DEFAULT_GROUP",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println(success)

	serviceList, err := client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		PageNo:   1,
		PageSize: 10,
	})
	if err != nil {
		panic(err)
	}
	for _, svc := range serviceList.Doms {
		fmt.Printf("Service: %s\n", svc)
		instances, err := client.SelectAllInstances(vo.SelectAllInstancesParam{
			ServiceName: svc,
		})
		if err != nil {
			panic(err)
		}
		bytes, _ := json.MarshalIndent(instances, "", " ")
		fmt.Println(string(bytes))
		for _, instance := range instances {
			fmt.Printf("Service instance: %s:%d\n", instance.Ip, instance.Port)
		}
	}
}
