package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hudl/fargo"
)

func main() {
	httpAddr := "http://127.0.0.1:8761/eureka"
	eurekaClient := fargo.NewConn(httpAddr)

	apps, err := eurekaClient.GetApps()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	fmt.Println(len(apps))

	bytes, _ := json.Marshal(apps)
	fmt.Println(string(bytes))
}
