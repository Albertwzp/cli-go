package main

import (
	"context"
	"fmt"

	kconfig "github.com/Albertwzp/cli-go/config"
	"github.com/Albertwzp/cli-go/tools"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

func main() {
	config, err := kconfig.GetK8sConfig()
	//参考path ： /api/v1/namespace/{namespace}/pods
	config.APIPath = "api"
	//pod的group是空字符串
	/*
	   // GroupName is the group name use in this package
	   const GroupName = ""
	   // SchemeGroupVersion is group version used to register these objects
	   var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
	*/

	config.GroupVersion = &corev1.SchemeGroupVersion
	//指定序列化工具
	config.NegotiatedSerializer = scheme.Codecs

	//根据配置信息构建restClient示例
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		fmt.Println("init restClient failed ! err: ", err)
		panic(err.Error())
	}

	//保存pod结果的数据结构实例
	result := &corev1.PodList{}

	//指定namespace
	namespace := ""

	//设置请求参数，然后发起请求

	//GET请求
	err = restClient.Get().
		//指定namespace，参考path: /api/v1/namespace/{namespace}/pods
		Namespace(namespace).
		// 查找多个pod，参考path: /api/v1/namespace/{namespace}/pods
		Resource("pods").
		//指定大小限制和序列化工具
		// 限制指定返回结果条目为100条
		//VersionedParams(&metav1.ListOptions{Limit:100}, scheme.ParameterCodec).
		//使用字段选择器选择只返回metadata.name为coredns-64dc4c69b-v8t4m的pod
		//VersionedParams(&metav1.ListOptions{Limit:100,FieldSelector:"metadata.name=coredns-64dc4c69b-v8t4m"}, scheme.ParameterCodec).
		//使用字段选择器选择只返回状态（status.phase）为Running的pod
		//VersionedParams(&metav1.ListOptions{Limit:100,FieldSelector:"status.phase=Running"}, scheme.ParameterCodec).
		//使用标签选择器选择标签k8s-app=kube-dns的pod
		//VersionedParams(&metav1.ListOptions{Limit:100,LabelSelector:"k8s-app=kube-dns"}, scheme.ParameterCodec).
		//请求
		Do(context.TODO()).
		//将结果存入result
		Into(result)

	if err != nil {
		fmt.Println("RESTClient Get failed ! err: ", err)
		panic(err.Error())
	}
	data := tools.FormatData(result)
	tools.PrintTab(data)
}
