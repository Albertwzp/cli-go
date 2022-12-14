package tools

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/Albertwzp/cli-go/config"
	vault "github.com/hashicorp/vault/api"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Conn struct {
	K8sC   []*kubernetes.Clientset
	VaultC *vault.Client
	Env    string
}

func CreateConn(env string, vi [2]string, ci [][2]string) *Conn {
	vclient := VaultC(vi[0], vi[1])
	var clsS []*kubernetes.Clientset
	for _, i := range ci {
		info := map[string]string{"env": env, "apiServer": i[0], "token": i[1]}
		cls := config.NewTokenConfig(info)
		clsS = append(clsS, cls)
	}
	return &Conn{
		K8sC:   clsS,
		VaultC: vclient,
		Env:    env,
	}
}

func FormatYaml(cmlist *v1.ConfigMapList, key string) map[string]string {
	var kv map[string]string
	for _, cm := range cmlist.Items {
		if strings.Contains(cm.Name, key+"-configmap") {
			data := cm.Data["application.yaml"]
			//yaml.Unmarshal()
			kv, _ = ConvertYaml(data)
		}
	}
	return kv
}
func (c *Conn) GetVt(app string) map[string]string {
	kv := VaultKv(c.VaultC, c.Env, app)
	return kv
}

func (c *Conn) Union(app string) map[string]string {
	cc := c.GetCm(app)
	//cv := make(map[string]string)
	cv := c.GetVt(app)
	return KvUnion(cc, cv)
}

func (c *Conn) GetPod() [][]string {
	var data [][]string
	for _, k := range c.K8sC {
		podClient := k.CoreV1().Pods("")
		resultC, err := podClient.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Println("podclient get pods failed! err: ", err)
			panic(err.Error())
		}
		rs := FormatPod(resultC)
		for _, x := range rs {
			tmp := x
			data = append(data, tmp)
		}
	}
	sort.Sort(sortable(data))
	return data
}
