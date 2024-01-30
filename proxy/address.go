package proxy

import "strings"

func isTailscaleHost(host string) bool {
	host = strings.ToLower(host)

	if strings.HasSuffix(host, ".ts.net") {
		return true
	}

	if strings.HasPrefix(host, "100.") || strings.HasPrefix(host, "fd7a:115c:a1e0:") {
		return true
	}

	return false
}
