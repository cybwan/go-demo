package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hudl/fargo"
)

func main() {
	httpAddr := "http://127.0.0.1:8761/eureka"
	//httpAddr := "http://192.168.10.47:8760/eureka"
	eurekaClient := fargo.NewConn(httpAddr)
	{
		apps, err := eurekaClient.GetApps()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}

		bytes, _ := json.MarshalIndent(apps, "", " ")
		fmt.Println(string(bytes))
	}
	//{
	//	app, err := eurekaClient.GetApp("httpbin")
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		os.Exit(-1)
	//	}
	//
	//	bytes, _ := json.MarshalIndent(app, "", " ")
	//	fmt.Println(string(bytes))
	//}
}
