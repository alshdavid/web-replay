package internal_serve

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/alshdavid/web-replay/src/platform/har"
)

type ServerSetting struct {
	OriginalHost string
	NewHost      string
	Entries      map[string][]har.Entry
	Domain       string
	Port         int
}

func GenerateSettings(serverPort int, serverDomain string, harEntries []har.Entry) []ServerSetting {
	serverSettings := []ServerSetting{}

	for _, entry := range harEntries {
		u, _ := url.Parse(entry.Request.Url)
		var setting *ServerSetting

		for _, v := range serverSettings {
			if v.OriginalHost == u.Host {
				setting = &v
				break
			}
		}
		if setting == nil {
			setting = &ServerSetting{
				OriginalHost: u.Host,
				NewHost:      fmt.Sprintf("https://%s:%d", serverDomain, serverPort),
				Domain:       serverDomain,
				Port:         serverPort,
				Entries:      map[string][]har.Entry{},
			}
			serverPort += 1
			serverSettings = append(serverSettings, *setting)
		}

		url := u.String()
		url = strings.TrimPrefix(url, u.Scheme)
		url = strings.TrimPrefix(url, "://")
		url = strings.TrimPrefix(url, u.Host)

		cacheKey := fmt.Sprintf("%s:%s", entry.Request.Method, url)
		_, ok := setting.Entries[cacheKey]
		if !ok {
			setting.Entries[cacheKey] = []har.Entry{}
		}
		setting.Entries[cacheKey] = append(setting.Entries[cacheKey], entry)
	}

	return serverSettings
}
