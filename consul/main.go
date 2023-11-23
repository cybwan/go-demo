package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
)

func main() {
	c := api.DefaultConfig()
	c.Address = "127.0.0.1:8500"

	if consulClient, err := api.NewClient(c); err != nil {
		fmt.Println(err.Error())
	} else {
		ctx := context.Background()
		opts := (&api.QueryOptions{
			AllowStale: true,
			WaitIndex:  1,
			WaitTime:   5 * time.Second,
		}).WithContext(ctx)
		//if serviceMap, _, err := consulClient.Catalog().Services(opts); err == nil {
		//	for svc := range serviceMap {
		//		if serviceEntries, _, err := consulClient.Health().Service(svc, "", false, nil); err == nil {
		//			for index, serviceEntry := range serviceEntries {
		//				bytes, _ := json.MarshalIndent(serviceEntry, "", " ")
		//				pwd, _ := os.Getwd()
		//				os.WriteFile(fmt.Sprintf("%s/data/%s.%d.json", pwd, svc, index), bytes, os.ModePerm)
		//			}
		//		}
		//	}
		//}

		//opts := &api.QueryOptions{
		//	AllowStale: true,
		//	WaitIndex:  1,
		//	WaitTime:   1 * time.Minute,
		//	Filter:     fmt.Sprintf("\"%s\" in Tags", s.ConsulK8STag),
		//}

		services, _, err := consulClient.Catalog().NodeServiceList("consul-578ccd7c59-ntw7g", opts)
		if err != nil {
			fmt.Println(err.Error())
		}
		bytes, _ := json.MarshalIndent(services, "", " ")
		fmt.Println(string(bytes))
	}
}
