package main

import (
	"fmt"

	config "github.com/Albertwzp/cli-go/Config"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

func main() {
	config, err := config.GetK8sConfig()
	//实例化一个discoveryClient对象
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		fmt.Println("init discoveryClient failed ! err: ", err)
		panic(err.Error())
	}

	//获取所有分组的资源数据
	APIGroup, APIResourceListSlice, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		fmt.Println("get apigroup resource failed! err: ", err)
		panic(err.Error())
	}
	//先看Group信息
	fmt.Printf("APIGroup : \n\n %v\n\n\n\n", APIGroup)

	// APIResourceListSlice是个切片，里面的每个元素代表一个GroupVersion及其资源
	for _, singleAPIResourceList := range APIResourceListSlice {

		// GroupVersion是个字符串，例如"apps/v1"
		groupVerionStr := singleAPIResourceList.GroupVersion

		// ParseGroupVersion方法将字符串转成数据结构
		gv, err := schema.ParseGroupVersion(groupVerionStr)

		if err != nil {
			panic(err.Error())
		}

		fmt.Println("*****************************************************************")
		fmt.Printf("GV string [%v]\nGV struct [%#v]\nresources :\n\n", groupVerionStr, gv)

		// APIResources字段是个切片，里面是当前GroupVersion下的所有资源
		for _, singleAPIResource := range singleAPIResourceList.APIResources {
			fmt.Printf("%v\n", singleAPIResource.Name)
		}
	}
}
