package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	consulclient *api.Client
	kubeclient   *kubernetes.Clientset
)

func getConsulClient() *api.Client {
	if consulclient == nil {
		c := api.DefaultConfig()
		c.Address = "127.0.0.1:8500"
		consulclient, _ = api.NewClient(c)
	}
	return consulclient
}

func getKubeClient() *kubernetes.Clientset {
	if kubeclient == nil {
		// Initialize kube config and client
		kubeConfigFile := "/Users/baili/.kube/config"
		kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigFile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
		kubeclient = kubernetes.NewForConfigOrDie(kubeConfig)
	}
	return kubeclient
}
