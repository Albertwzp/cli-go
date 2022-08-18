package main

import (
	"context"
	"fmt"

	kconfig "github.com/Albertwzp/cli-go/config"
	"github.com/Albertwzp/cli-go/tools"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func main() {
	config, err := kconfig.GetK8sConfig()
	//实例化一个clientset对象
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("init clientset failed ! err: ", err)
		panic(err.Error())
	}

	//获取podClient客户端,corev1.NamespaceAll 为空字符串，实际如果为空字符串，那么拿到的是所有名称空间的pod资源
	//podClient := clientset.CoreV1().Pods(corev1.NamespaceAll)
	podClient := clientset.CoreV1().Pods("")

	//使用podclient客户端，列出名称空间内所有pod资源
	result, err := podClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("podclient get pods failed! err: ", err)
		panic(err.Error())
	}
	data := tools.FormatData(result)
	tools.PrintTab(data)
}
