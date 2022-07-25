package execute

import (
	"net/http"
	"strings"
)

func isJSONResponse(r *http.Response) bool {
	contentType := r.Header.Get("Content-Type")
	return contentType == "application/json" || strings.HasPrefix(contentType, "application/json;")
}
