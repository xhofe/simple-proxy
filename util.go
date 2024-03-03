package main

import (
	"net/url"
	"strings"
)

func getKeysByPrefix(query url.Values, prefix string) (keys []string) {
	for k, _ := range query {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	return
}
