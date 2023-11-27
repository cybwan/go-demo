package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/cybwan/go-demo/catalog"
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
			//for svc := range serviceMap {
			//	if serviceEntries, _, err := consulClient.Health().CatalogService(svc, "", false, nil); err == nil {
			//		for index, serviceEntry := range serviceEntries {
			//			bytes, _ := json.MarshalIndent(serviceEntry, "", " ")
			//			pwd, _ := os.Getwd()
			//			os.WriteFile(fmt.Sprintf("%s/data/%s.%d.json", pwd, svc, index), bytes, os.ModePerm)
			//		}
			//	}
			//}

			//bytes, _ := json.MarshalIndent(servicesMap, "", " ")
			//fmt.Println(string(bytes))

			for svc, svcTags := range servicesMap {
				if strings.EqualFold(svc, "consul") {
					continue
				}
				for idx, svcTag := range svcTags {
					tags := catalog.ParseTags(svcTag)
					for _, tag := range tags {
						k, v := catalog.ParseTag(tag)
						fmt.Printf("service=%s index=%d k=%s v=%s\n", svc, idx, k, v)
					}
					//bytes, _ := json.MarshalIndent(tags, "", " ")
					//fmt.Println(string(bytes))
				}
			}
		}

		//opts := &api.QueryOptions{
		//	AllowStale: true,
		//	WaitIndex:  1,
		//	WaitTime:   1 * time.Minute,
		//	Filter:     fmt.Sprintf("\"%s\" in Tags", s.ConsulK8STag),
		//}

		//services, _, err := consulClient.Catalog().NodeServiceList("consul-578ccd7c59-ntw7g", opts)
		//if err != nil {
		//	fmt.Println(err.Error())
		//}
		//bytes, _ := json.MarshalIndent(services, "", " ")
		//fmt.Println(string(bytes))
	}
}
