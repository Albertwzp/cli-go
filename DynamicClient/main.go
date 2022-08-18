package main

import (
	"context"
	"fmt"

	"github.com/Albertwzp/cli-go/config"
	"github.com/Albertwzp/cli-go/tools"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func main() {
	configk, err := config.GetK8sConfig()
	//实例化一个DynamicClient对象
	dynamicClient, err := dynamic.NewForConfig(configk)
	if err != nil {
		fmt.Println("init dynamicClient failed ! err: ", err)
		panic(err.Error())
	}

	//dynamicClient的唯一关联方法所需的入参，GVR
	gvr := schema.GroupVersionResource{Version: "v1", Resource: "pods"}

	//使用dynamicClient的查询列表方法，查询指定namespace下的所有pod
	//注意此方法返回的数据结构类型是UnstructuredList
	unstructObjList, err := dynamicClient.
		//Resource是dynamicClient唯一的一个方法，参数为gvr
		Resource(gvr).
		//指定查询的namespace
		Namespace("").
		//以list列表的方式查询
		List(context.TODO(), metav1.ListOptions{Limit: 100})

	if err != nil {
		fmt.Println("dynamicClient list pods failed ! err :", err)
		panic(err.Error())
	}

	//实例化一个PodList数据结构，用于接收从unstructObjList转换后的结果
	result := &corev1.PodList{}

	//转换
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObjList.UnstructuredContent(), result)
	if err != nil {
		fmt.Println("unstructured failed ! err: ", err)
		panic(err.Error())
	}
	data := tools.FormatData(result)
	tools.PrintTab(data)
}
