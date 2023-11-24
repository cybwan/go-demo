package catalog

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeclient *kubernetes.Clientset
)

func GetKubeClient() *kubernetes.Clientset {
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
