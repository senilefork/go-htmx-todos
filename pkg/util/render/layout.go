package render

import (
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"go-htmx-jobs/pkg/util"
)

// TargetFunc wraps page content in one of the layout shells (full document,
// page layout only, ...). TemplateTargets maps a DOM id (or "index") to one.
type TargetFunc func(content templ.Component) templ.ComponentFunc
type TemplateTargets map[string]TargetFunc

// Layout decides how much HTML to send back:
//
//   - Not an HTMX request (first page load) -> wrap in "index": the full document.
//   - HTMX request with HX-Target           -> wrap in just that target: a fragment.
//   - Anything else                         -> render the bare page component.
//
// The same handler therefore serves both full page loads and HTMX partials.
// templateTargets must contain an "index" entry.
func Layout(ctx context.Context, req *http.Request, w io.Writer, templateTargets TemplateTargets, page templ.Component) error {
	hxRequest := util.IsHxRequest(req)
	target := req.Header.Get("HX-Target")

	if !hxRequest {
		if tt, ok := templateTargets["index"]; ok {
			return tt(page).Render(ctx, w)
		}
		panic("no index template")
	}

	if tt, ok := templateTargets[target]; ok {
		return tt(page).Render(ctx, w)
	}

	return page.Render(ctx, w)
}
