package internal_serve

import "fmt"

func GenerateDomainMap(serverSettings []ServerSetting) map[string]string {
	domainMap := map[string]string{}
	for _, v := range serverSettings {
		domainMap[v.OriginalHost] = fmt.Sprintf("https://%s:%d", v.Domain, v.Port)
	}
	return domainMap
}
