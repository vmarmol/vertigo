package static

import (
	"fmt"
	"net/http"
	"net/url"
)

const StaticResource = "/static/"

var staticFiles = map[string]string{
	"index.html": indexHtml,
	"styles.css": stylesCss,
	"scripts.js": scriptsJs,
}

func HandleRequest(w http.ResponseWriter, u *url.URL) error {
	if len(u.Path) <= len(StaticResource) {
		return fmt.Errorf("unknown static resource %q", u.Path)
	}

	// Get the static content if it exists.
	resource := u.Path[len(StaticResource):]
	content, ok := staticFiles[resource]
	if !ok {
		return fmt.Errorf("unknown static resource %q", resource)
	}

	_, err := w.Write([]byte(content))
	return err
}
