package tools

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

func VaultC(srv, token string) *vault.Client {
	config := vault.DefaultConfig()
	config.ConfigureTLS(&vault.TLSConfig{Insecure: true})
	config.Address = srv
	//config.ConfigureTLS(clientTLSConfig.InsecureSkipVerify)
	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	// Authenticate
	client.SetToken(token)
	return client
}

func VaultM(client *vault.Client, env string) map[string][]string {
	mounts, _ := client.Sys().ListMounts()
	apps := make(map[string][]string)
	for k, _ := range mounts {
		app := k[:len(k)-1]
		if strings.Contains(k, env) {
			if logic, _ := client.Logical().List(app); logic != nil {
				for _, v := range logic.Data["keys"].([]interface{}) {
					apps[app] = append(apps[app], v.(string))
				}
				/*apps[app] = logic.Data["keys"].([]string)
				if ks, ok := logic.Data["keys"].([]string); ok {
					fmt.Println(ks)
					apps[app] = ks
				}
				m, _ := client.Sys().ListPlugins(&vault.ListPluginsInput{})*/
			}
		}
	}
	//fmt.Println("%v", apps)
	return apps
}

func VaultS(client *vault.Client, env, app string) map[string][]string {
	appx := make(map[string][]string)
	mm := VaultM(client, env)
	for m, _ := range mm {
		if strings.Replace(strings.TrimPrefix(m, "app-"), "-"+env, "", 1) == app {
			appx[m] = mm[m]
		}
	}

	return appx
}

func VaultKv(client *vault.Client, env, app string) map[string]string {
	kv := make(map[string]string)
	ap := VaultS(client, env, app)
	for k, vs := range ap {
		for _, v := range vs {
			xy, _ := client.KVv1(k).Get(context.Background(), v)
			for x, y := range xy.Data {
				switch val := y.(type) {
				case json.Number:
					kv[x] = val.String()
				case string:
					kv[x] = val
				default:
					kv[x] = ""
				}
			}
		}
	}
	return kv
}
