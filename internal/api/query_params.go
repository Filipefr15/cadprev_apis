package api

import (
	"net/url"
)

// BuildQueryString monta a query string a partir de um map de parâmetros
func BuildQueryString(params map[string]string) string {
	q := url.Values{}
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}
	return q.Encode()
}
