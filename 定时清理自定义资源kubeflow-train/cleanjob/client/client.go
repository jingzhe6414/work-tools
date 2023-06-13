package client

import (
	"fmt"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

func NewClient() *dynamic.DynamicClient {
	kubeConfigPath := ""
	if home := homedir.HomeDir(); home != "" {
		kubeConfigPath = filepath.Join(home, ".kube", "config")
	}
	fileExist, err := PathExists(kubeConfigPath)
	if err != nil {
		fmt.Println("默认config不存在")
	}
	if fileExist {
		fmt.Println("从默认路径读取")

		config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			panic(err.Error())
		}
		// creates the clientset
		dynamicClient, err := dynamic.NewForConfig(config)
		//clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		return dynamicClient
	} else {
		fmt.Println("读取内置")
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		// creates the clientset
		dynamicClient, err := dynamic.NewForConfig(config)
		//clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		return dynamicClient
	}

}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
