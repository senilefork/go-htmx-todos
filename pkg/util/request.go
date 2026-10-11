package util

import "net/http"

// IsHxRequest reports whether the request came from HTMX (it always sends
// the HX-Request header). This is how the server knows whether to return a
// full page or just a fragment.
func IsHxRequest(req *http.Request) bool {
	return req.Header.Get("HX-Request") == "true"
}

// HandleRedirect sends HTMX's client-side redirect header when the request
// came from HTMX, and a normal 302 otherwise. Same UX, different mechanism.
func HandleRedirect(w http.ResponseWriter, req *http.Request, url string) {
	if IsHxRequest(req) {
		w.Header().Set("HX-Redirect", url)
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, req, url, http.StatusFound)
}
