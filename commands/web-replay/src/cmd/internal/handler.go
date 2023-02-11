package internal_serve

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/alshdavid/web-replay/src/platform/extras"
	"github.com/andybalholm/brotli"
	"golang.org/x/term"
)

type CachedBody struct {
	Headers         http.Header
	Body            [][]byte
	PayloadHash     string
	ResponseCode    int
	ContentEncoding string
	ExtraLatency    time.Duration
	StreamLatency   time.Duration
}

func Handler(
	setting ServerSetting,
	domainMap map[string]string,
	blockXhrRequestTime time.Duration,
	servicePatches []Patch,
	logger ILogger,
) http.HandlerFunc {
	var cache = map[string]CachedBody{}

	// Pre-process requests applying patches and compression
	for _, routeEntries := range setting.Entries {
		for _, entry := range routeEntries {
			var responseBodyChunks [][]byte
			var responseBody []byte
			var headers http.Header = http.Header{}
			var extraLatency time.Duration
			var streamLatency time.Duration
			var mimeType = entry.Response.Content.MimeType
			var textEncoding = entry.Response.Content.Encoding
			var contentText = entry.Response.Content.Text
			var cachedRoutePath = entry.Request.Url
			var contentEncoding string

			cachedRoutePath = strings.TrimPrefix(cachedRoutePath, "https://")
			cachedRoutePath = strings.TrimPrefix(cachedRoutePath, "http://")
			cacheKey := fmt.Sprintf("%s:%s", entry.Request.Method, cachedRoutePath)

			if entry.Request.PostData.Text != "" {
				payloadHash := extras.GetMD5Hash(entry.Request.PostData.Text)
				cacheKey = fmt.Sprintf("%s:%s", cacheKey, payloadHash)
			}

			for _, h := range entry.Response.Headers {
				if !isHeaderAllowed(h.Name) {
					continue
				}
				headers.Set(h.Name, transformDomains(h.Value, domainMap))
			}

			if isWriteableMime(mimeType) {
				replacedText := contentText

				if textEncoding == "base64" {
					r, _ := base64.StdEncoding.DecodeString(replacedText)
					replacedText = string(r)
				}
				replacedText = transformDomains(replacedText, domainMap)
				responseBody = []byte(replacedText)
			}

			if !isWriteableMime(mimeType) && textEncoding != "base64" {
				responseBody = []byte(contentText)
			}

			if textEncoding == "base64" {
				r, _ := base64.StdEncoding.DecodeString(contentText)
				responseBody = r
			}

			// Apply patches
			for i := len(servicePatches) - 1; i >= 0; i-- {
				patch := servicePatches[i]
				mimeMatch := patch.Match.MimeTypes == nil
				originMatch := patch.Match.Origins == nil
				urlPathMatch := patch.Match.UrlPaths == nil
				requestTypeMatch := patch.Match.RequestTypes == nil

				if len(patch.Match.MimeTypes) != 0 &&
					extras.SliceContains(patch.Match.MimeTypes, mimeType) {
					mimeMatch = true
				}

				if len(patch.Match.Origins) != 0 &&
					extras.SliceContains(patch.Match.Origins, setting.OriginalHost) {
					originMatch = true
				}

				if len(patch.Match.RequestTypes) != 0 &&
					extras.SliceContains(patch.Match.RequestTypes, entry.ResourceType) {
					requestTypeMatch = true
				}

				if len(patch.Match.UrlPaths) != 0 {
					for _, pattern := range patch.Match.UrlPaths {
						ok, _ := filepath.Match(pattern, cachedRoutePath)
						if ok {
							urlPathMatch = true
							break
						}
					}
				}

				if !(mimeMatch && originMatch && urlPathMatch && requestTypeMatch) {
					continue
				}

				if patch.AddLatency != nil {
					extraLatency = time.Duration((*patch.AddLatency) * int64(time.Millisecond))
				}

				if patch.ReplaceText != nil {
					text := string(responseBody)
					text = strings.ReplaceAll(text, patch.ReplaceText.FindString, patch.ReplaceText.ReplaceWith)
					responseBody = []byte(text)
				}
			}

			for i := len(servicePatches) - 1; i >= 0; i-- {
				patch := servicePatches[i]
				if patch.StreamLatency != nil && strings.Contains(string(responseBody), patch.StreamLatency.FindString) {
					streamLatency = time.Duration(*patch.StreamLatency.Latency) * time.Millisecond
					text := string(responseBody)
					splits := strings.Split(text, patch.StreamLatency.FindString)
					if len(splits) > 1 {
						splits[0] += patch.StreamLatency.FindString
					}
					for _, text := range splits {
						responseBodyChunks = append(responseBodyChunks, []byte(text))
					}
				}
			}

			if len(responseBodyChunks) == 0 && entry.ResourceType != "fetch" {
				contentEncoding = "br"

				var b bytes.Buffer
				bw := brotli.NewWriter(&b)
				bw.Write(responseBody)
				bw.Close()

				responseBodyChunks = append(responseBodyChunks, b.Bytes())
			}

			if len(responseBodyChunks) == 0 && entry.ResourceType == "fetch" {
				responseBodyChunks = append(responseBodyChunks, []byte(responseBody))
			}

			cache[cacheKey] = CachedBody{
				Body:            responseBodyChunks,
				Headers:         headers,
				ContentEncoding: contentEncoding,
				ExtraLatency:    extraLatency,
				StreamLatency:   streamLatency,
				ResponseCode:    entry.Response.Status,
			}
		}
	}

	// Runs on each request
	return func(w http.ResponseWriter, r *http.Request) {
		width, _, _ := term.GetSize(0)
		routePath := fmt.Sprintf("%s%s", setting.OriginalHost, r.URL)
		logOutput := fmt.Sprintf("%-7s https://localhost:%d => %s", r.Method, setting.Port, r.URL)

		logger.Printf("%s\n", extras.TruncateString(logOutput, width))

		cacheKey := fmt.Sprintf("%s:%s", r.Method, routePath)

		defer r.Body.Close()
		bodyBytes, _ := io.ReadAll(r.Body)
		bodyString := string(bodyBytes)

		if bodyString != "" {
			payloadHash := extras.GetMD5Hash(bodyString)
			cacheKey = fmt.Sprintf("%s:%s", cacheKey, payloadHash)
		}

		cachedResponse, ok := cache[cacheKey]

		if !ok {
			w.WriteHeader(500)
			w.Write([]byte("No entry cached"))
			return
		}

		for key, values := range cachedResponse.Headers {
			for _, value := range values {
				w.Header().Set(key, value)
			}
		}

		if cachedResponse.ExtraLatency != 0 {
			time.Sleep(cachedResponse.ExtraLatency)
		}

		if cachedResponse.ContentEncoding != "" {
			w.Header().Set("Content-Encoding", cachedResponse.ContentEncoding)
		}

		w.WriteHeader(cachedResponse.ResponseCode)

		if len(cachedResponse.Body) == 1 {
			w.Write(cachedResponse.Body[0])
			return
		}

		for i, chunk := range cachedResponse.Body {
			w.Write(chunk)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			if i != len(cachedResponse.Body) {
				time.Sleep(cachedResponse.StreamLatency)
			}
		}

	}
}
