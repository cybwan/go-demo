package main

import (
	"github.com/hudl/fargo"
)

func main() {
	httpAddr := "http://127.0.0.1:8761/eureka"
	eurekaClient := fargo.NewConn(httpAddr)
	bookstoreIns := fargo.Instance{
		HostName:         "192.168.226.21",
		InstanceId:       "bookstore-192.168.226.21",
		App:              "BOOKSTORE",
		Port:             9091,
		IPAddr:           "192.168.226.21",
		VipAddress:       "bookstore",
		SecureVipAddress: "bookstore",
		DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
		Status:           fargo.UP,
	}
	eurekaClient.RegisterInstance(&bookstoreIns)

	bookbuyerIns := fargo.Instance{
		HostName:         "192.168.226.22",
		InstanceId:       "bookbuyer-192.168.226.22",
		App:              "BOOKBUYER",
		Port:             9092,
		IPAddr:           "192.168.226.22",
		VipAddress:       "bookbuyer",
		SecureVipAddress: "bookbuyer",
		DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
		Status:           fargo.UP,
	}
	eurekaClient.RegisterInstance(&bookbuyerIns)
}
