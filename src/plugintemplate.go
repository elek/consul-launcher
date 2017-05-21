package consullauncher

import (
	"html/template"
	"github.com/Masterminds/sprig"
	"fmt"
	"bytes"
	"github.com/hashicorp/consul/api"
	"strconv"
)

func templatePostIteration(configFiles []Entry, dest string) {

}

var templatePlugin = Plugin{
	PostIteration: func([] Entry, string) {},
	CheckActivation:func(flag uint64) bool {
		return flag | 2 > 0
	},
	ProcessContent:func(content []byte, consul *api.Client, supervisor chan bool) []byte {
		funcMap := sprig.FuncMap()
		funcMap["service"] = func(service string) ([]*api.CatalogService, error) {
			entry, queryMeta, err := consul.Catalog().Service(service, "", nil)
			go listenToChanges(consul, service, queryMeta.LastIndex, supervisor)
			if err != nil {
				return nil, err
			} else {
				return entry, nil
			}
		}
		funcMap["hosts"] = func(separator string, services []*api.CatalogService) string {
			result := ""
			for ix, service := range services {
				if (ix > 0 && len(separator) > 0) {
					result += separator
				}
				result += service.Address

			}
			return result

		}
		funcMap["hostPorts"] = func(separator string, services []*api.CatalogService) string {
			result := ""
			for ix, service := range services {
				if (ix > 0 && len(separator) > 0) {
					result += separator
				}
				result += (service.Address + ":" + strconv.Itoa(service.ServicePort))

			}
			return result

		}
		template, err := template.New("template").Funcs(funcMap).Parse(string(content))
		if (err != nil) {
			fmt.Println("Error on parsing template")
			return content
		}
		var result bytes.Buffer
		template.Execute(&result, make(map[string]string))
		return result.Bytes()
	},
}

var services = make(map[string]bool)

func listenToChanges(consul *api.Client, service string, lastIndex uint64, supervisor chan bool) {
	if _, ok := services[service]; !ok {
		services[service] = true
		println("Start to listen on changes of service " + service)
		currentIndex := lastIndex
		for {
			queryOptions := api.QueryOptions{
				WaitIndex: currentIndex,
			}
			_, queryMeta, err := consul.Catalog().Service(service, "", &queryOptions)
			if (err == nil && queryMeta.LastIndex > currentIndex) {
				fmt.Printf("Service %s has been changed from %d to %d\n", service, currentIndex, queryMeta.LastIndex)
				currentIndex = queryMeta.LastIndex
				supervisor <- true
			}
		}
	}
}