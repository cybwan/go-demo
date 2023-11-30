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
		if servicesMap, _, err := consulClient.Catalog().Services(opts); err == nil {
			for svc := range servicesMap {
				if serviceEntries, _, err := consulClient.Health().Service(svc, "", false, nil); err == nil {
					for index, serviceEntry := range serviceEntries {
						bytes, _ := json.MarshalIndent(serviceEntry, "", " ")
						pwd, _ := os.Getwd()
						os.WriteFile(fmt.Sprintf("%s/data/%s.%d.json", pwd, svc, index), bytes, os.ModePerm)
					}
				}
			}

			bytes, _ := json.MarshalIndent(servicesMap, "", " ")
			fmt.Println(string(bytes))
		}
	}
}
