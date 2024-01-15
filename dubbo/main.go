package main

import (
	"encoding/json"
	"fmt"

	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"

	"github.com/cybwan/go-demo/dubbo/zookeeper"
)

func main() {
	url, _ := common.NewURL("dubbo://127.0.0.1:2181",
		common.WithParamsValue(constant.ClientNameKey, "zk-client"))
	sd, err := zookeeper.NewZookeeperServiceDiscovery(url)

	if err != nil {
		panic(err)
	}

	//services := sd.GetServices()
	//
	//for svc, _ := range services.Items {
	//
	//}

	instances := sd.GetInstances("dubbo_registry_zookeeper_server")
	for _, v := range instances {
		bytes, _ := json.MarshalIndent(v, "", " ")
		fmt.Println(string(bytes))
	}

	err = sd.Destroy()

	if err != nil {
		panic(err)
	}
}
