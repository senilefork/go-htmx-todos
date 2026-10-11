package common

import (
	"net/http"

	"go-htmx-jobs/internal/templates/ui/partials"
	"go-htmx-jobs/pkg/util"
)

// HandleError adapts the response to how the request arrived
// (trimmed version of the starter's internal/common/reqerror.go):
//
//   - HTMX request -> HTML toast fragment + hx-toast header + real status;
//     the global listener in layouts.Index injects it into #global-toast.
//   - Normal request -> plain http.Error text.
func HandleError(req *http.Request, w http.ResponseWriter, status int, message string) {
	if util.IsHxRequest(req) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("hx-toast", "true")
		w.WriteHeader(status)
		_ = partials.Toast(message).Render(req.Context(), w)
		return
	}
	http.Error(w, message, status)
}
