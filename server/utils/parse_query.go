package utils

import (
	"net/http"
	"strings"
)

func ParseQueryList(r *http.Request, param string, sep string) []string {
	var res []string
	rawStrs := r.URL.Query()[param]
	for _, rawStr := range rawStrs {
		for part := range strings.SplitSeq(rawStr, sep) {
			part = strings.TrimSpace(part)
			if part != "" {
				res = append(res, part)
			}
		}
	}
	return res
}
