package public

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"go-htmx-jobs/internal/common"
	appctx "go-htmx-jobs/internal/context"
	"go-htmx-jobs/internal/services"
	"go-htmx-jobs/internal/templates/constants"
	"go-htmx-jobs/internal/templates/templatetargets"
	"go-htmx-jobs/internal/templates/ui/pages"
	"go-htmx-jobs/internal/templates/ui/tables"
	"go-htmx-jobs/pkg/util"
	"go-htmx-jobs/pkg/util/render"
)

// todos is one controller wrapping jsonplaceholder's /todos endpoints,
// same shape as the starter's public/home.go or private/management.go.
type todos struct {
	*appctx.AppContext
}

func NewTodos(ctx *appctx.AppContext) *todos {
	return &todos{AppContext: ctx}
}

func (c *todos) Path() string {
	return "/" // mounted last: /public is registered before it
}

func (c *todos) PublicRoute(r chi.Router) {
	r.Get("/", c.redirect)
	r.Get("/todos", c.list)
	r.Get("/todos/{id}", c.detail)
}

func (c *todos) redirect(w http.ResponseWriter, req *http.Request) {
	util.HandleRedirect(w, req, "/todos")
}

// list handles the page and the search: the search input sends
// HX-Target: todos-table, so only the table component is returned.
func (c *todos) list(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := appctx.GetLogger(ctx, "todos.list")

	query := req.URL.Query().Get("q")
	todos, err := c.SM().TodosService().List(ctx, query)
	if err != nil {
		lgr.Error("failed to list todos", "error", err)
		common.HandleError(req, w, http.StatusBadGateway, "Failed to load todos")
		return
	}
	lgr.Info("listed todos", "count", len(todos), "query", query)

	var content templ.Component = pages.Todos(todos, query)
	if req.Header.Get("HX-Target") == constants.TodosTableID {
		content = tables.TodosTable(todos)
	}

	if err := render.Layout(ctx, req, w, templatetargets.Main, content); err != nil {
		lgr.Error("render failed", "error", err)
	}
}

// detail serves /todos/{id} by calling the upstream detail endpoint.
// A missing id upstream is a 404; any other upstream failure is a 502.
func (c *todos) detail(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := appctx.GetLogger(ctx, "todos.detail")

	id, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		common.HandleError(req, w, http.StatusNotFound, "Todo not found")
		return
	}

	todo, err := c.SM().TodosService().GetByID(ctx, id)
	if errors.Is(err, services.ErrNotFound) {
		lgr.Info("todo not found", "id", id)
		common.HandleError(req, w, http.StatusNotFound, "Todo not found: "+strconv.Itoa(id))
		return
	}
	if err != nil {
		lgr.Error("failed to fetch todo", "id", id, "error", err)
		common.HandleError(req, w, http.StatusBadGateway, "Failed to load todo")
		return
	}

	if err := render.Layout(ctx, req, w, templatetargets.Main, pages.TodoDetail(todo)); err != nil {
		lgr.Error("render failed", "error", err)
	}
}
