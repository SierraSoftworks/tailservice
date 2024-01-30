package proxy

import "strings"

func isDebugLog(s string, args ...any) bool {
	return !strings.HasPrefix(s, "To start this tsnet server, restart with TS_AUTHKEY set")
}
