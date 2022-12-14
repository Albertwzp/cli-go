package tools

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/scylladb/termtables"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
)

func Html(kv map[string]string) string {
	t := termtables.CreateTable()
	t.AddHeaders("Key", "Value")
	for k, v := range kv {
		t.AddRow(k, v)
	}
	t.SetModeHTML()
	return t.Render()
}

func PrintTab(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAMESPACE", "NAME", "READY", "STATUS", "RESTART", "AGE", "ELAPSE"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}

func FormatData(result *corev1.PodList) (data [][]string) {
	data = make([][]string, 0, len(result.Items))
	for i, _ := range result.Items {
		var count int
		for _, v := range result.Items[i].Status.ContainerStatuses {
			if v.Ready {
				count++
			}
		}
		ss := make(map[v1.PodConditionType]time.Time)
		for _, s := range result.Items[i].Status.Conditions {
			ss[s.Type] = s.LastTransitionTime.Time
		}

		var avg string
		subHoure := time.Now().Sub(result.Items[i].ObjectMeta.CreationTimestamp.Time).Hours()
		if subHoure < 24 {
			avg = fmt.Sprintf("%fh", subHoure)
		} else {
			hours := int(subHoure)
			days := hours / 24
			if (hours % 24) > 0 {
				days += 1
			}
			avg = fmt.Sprintf("%dd", days)
		}

		data = append(data, []string{result.Items[i].Namespace, result.Items[i].Name, fmt.Sprintf("%d/%d", count, len(result.Items[i].Status.ContainerStatuses)),
			string(result.Items[i].Status.Phase), strconv.Itoa(int(result.Items[i].Status.ContainerStatuses[0].RestartCount)), avg, fmt.Sprintf("%fs", ss[v1.PodReady].Sub(ss[v1.PodScheduled]).Seconds())})
	}
	return data
}

func ConvertYaml(data string) (kv map[string]string, err error) {
	cs := make(map[string]string)
	viper.SetConfigType("YAML")
	if err := viper.ReadConfig(bytes.NewBuffer([]byte(data))); err != nil {
		return nil, err
	}
	ks := viper.AllKeys()
	for _, k := range ks {
		cs[k] = viper.GetString(k)
	}
	return cs, nil
}

func KvUnion(cm, vault map[string]string) map[string]string {
	reHide := regexp.MustCompile(`password|secret|[Kk]ey`)
	reMongo := regexp.MustCompile(`:?\w+@`)
	for k, v := range vault {
		if reHide.MatchString(k) {
			cm[k] = "******"
		} else if strings.Contains(v, "mongodb://") {
			txt := reMongo.FindString(v)
			cm[k] = strings.ReplaceAll(v, txt[1:len(txt)-1], "***")
		} else {
			cm[k] = v
		}
	}
	return cm
}

func FormatPod(result *corev1.PodList) (data [][]string) {
	data = make([][]string, 0, len(result.Items))
	for i, _ := range result.Items {
		if strings.Contains(result.Items[i].Namespace, "http") && strings.Contains(result.Items[i].Name, "dep") {
			ss := make(map[v1.PodConditionType]time.Time)
			for _, s := range result.Items[i].Status.Conditions {
				if s.Status == "True" {
					ss[s.Type] = s.LastTransitionTime.Time
				} else {
					ss[s.Type] = time.Now()
				}
			}

			subS := ss[v1.PodReady].Sub(ss[v1.PodScheduled]).Seconds()
			var startup string
			if subS > 600 {
				startup = fmt.Sprintf("%s", "inf")
			} else if subS < 0 {
				startup = fmt.Sprintf("%s", "nan")
			} else {
				startup = fmt.Sprintf("%9.1fs", subS)
			}

			data = append(data, []string{result.Items[i].Namespace, result.Items[i].Name, startup, string(result.Items[i].Status.Phase)})
		}
	}
	return data
}
