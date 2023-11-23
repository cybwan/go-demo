package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/consul/api"
)

func main() {
	c := api.DefaultConfig()
	c.Address = "192.168.10.91:8500"

	if consulClient, err := api.NewClient(c); err != nil {
		fmt.Println(err.Error())
	} else {
		ctx := context.Background()
		opts := (&api.QueryOptions{
			AllowStale: true,
			WaitIndex:  1,
			WaitTime:   5 * time.Second,
		}).WithContext(ctx)
		if serviceMap, _, err := consulClient.Catalog().Services(opts); err == nil {
			for svc := range serviceMap {
				if serviceEntries, _, err := consulClient.Health().Service(svc, "", false, nil); err == nil {
					for index, serviceEntry := range serviceEntries {
						bytes, _ := json.MarshalIndent(serviceEntry, "", " ")
						pwd, _ := os.Getwd()
						os.WriteFile(fmt.Sprintf("%s/data/%s.%d.json", pwd, svc, index), bytes, os.ModePerm)
					}
				}
			}
		}
	}
}
