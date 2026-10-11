package public

import (
	"net/http"
	"os"
	"path"

	"go-htmx-jobs/internal/context"
	"github.com/go-chi/chi/v5"
)

// webPublic serves the /public directory (css, images, ...) — static
// serving is just another controller in this architecture.
type webPublic struct {
	*context.AppContext
}

func NewWebPublic(ctx *context.AppContext) *webPublic {
	return &webPublic{AppContext: ctx}
}

func (wp *webPublic) Path() string {
	return "/public"
}

func (wp *webPublic) PublicRoute(r chi.Router) {
	r.Handle("/*", wp.ServePublicDir())
}

func (wp *webPublic) ServePublicDir() http.Handler {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	handler := http.FileServer(http.Dir(path.Join(wd, wp.Path())))
	return http.StripPrefix(wp.Path(), handler)
}
