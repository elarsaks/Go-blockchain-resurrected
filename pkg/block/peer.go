package block

import "strings"

func peerEndpoint(neighbor, path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	neighbor = strings.TrimRight(neighbor, "/")
	if strings.HasPrefix(neighbor, "http://") || strings.HasPrefix(neighbor, "https://") {
		return neighbor + path
	}
	return "http://" + neighbor + path
}
