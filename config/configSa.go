package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ClusterInfo struct {
	InsecureSkipTlsVerify    bool   `yaml:"insecure-skip-tls-verify,omitempty"`
	CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty" `
	Server                   string
}

type Cluster struct {
	Cluster ClusterInfo `yaml:"cluster"`
	Name    string      `yaml:"name"`
}

type ClusterList []Cluster

type ContextInfo struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type Context struct {
	Context ContextInfo `yaml:"context"`
	Name    string      `yaml:"name"`
}

type ContextList []Context

type UserInfo struct {
	Token                 string `yaml:"token,omitempty"`
	ClientCertificateData string `yaml:"client-certificate-data,omitempty"`
	ClientKeyData         string `yaml:"client-key-data,omitempty"`
}
type User struct {
	Name string   `yaml:"name"`
	User UserInfo `yaml:"user"`
}

type UserList []User

type Config struct {
	ApiVersion     string      `yaml:"apiVersion"`
	Clusters       ClusterList `yaml:"clusters,flow"`
	Contexts       ContextList `yaml:"contexts,flow"`
	CurrentContext string      `yaml:"current-context"`
	Kind           string      `yaml:"kind"`
	Users          UserList    `yaml:"users,flow"`
}

func NewTokenConfig(EnvConfig map[string]string) *kubernetes.Clientset {

	var cf = Config{
		ApiVersion: "v1",
		Clusters: ClusterList{
			{
				Name: EnvConfig["env"],
				Cluster: ClusterInfo{
					InsecureSkipTlsVerify: true,
					Server:                EnvConfig["apiServer"],
				},
			},
		},
		Kind: "Config",
		Users: UserList{
			{
				Name: EnvConfig["env"],
				User: UserInfo{
					Token: EnvConfig["token"],
				},
			},
		},
		Contexts: ContextList{
			{
				Context: ContextInfo{
					Cluster: EnvConfig["env"],
					User:    EnvConfig["env"],
				},
				Name: EnvConfig["env"],
			},
		},
		CurrentContext: EnvConfig["env"],
	}
	configResult, err := yaml.Marshal(&cf)
	if err != nil {
		panic(err.Error())
	}

	file, err := ioutil.TempFile(".", "tmp")
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		file.Close()
		os.RemoveAll(file.Name())
	}()
	file.WriteString(string(configResult))
	config, err := clientcmd.BuildConfigFromFlags("", file.Name())
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("init clientset failed ! err: ", err)
		panic(err.Error())
	}
	return clientset
}
