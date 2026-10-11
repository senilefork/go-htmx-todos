package templatetargets

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"go-htmx-jobs/internal/templates/constants"
	"go-htmx-jobs/internal/templates/ui/layouts"
	"go-htmx-jobs/pkg/util/render"
)

// Main is the set of layout shells render.Layout can choose from.
//
// "index"     -> full HTML document (head, htmx script, nav) for first loads
// "page-layout" -> everything below the nav, when HX-Target is that id
var Main = render.TemplateTargets{
	"index": func(content templ.Component) templ.ComponentFunc {
		return func(ctx context.Context, w io.Writer) error {
			return layouts.Index().Render(templ.WithChildren(ctx, templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				return layouts.MainPageLayout().Render(templ.WithChildren(ctx, content), w)
			})), w)
		}
	},
	constants.PageLayoutID: func(content templ.Component) templ.ComponentFunc {
		return func(ctx context.Context, w io.Writer) error {
			return layouts.MainPageLayout().Render(templ.WithChildren(ctx, content), w)
		}
	},
}
