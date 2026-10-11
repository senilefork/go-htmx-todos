# go-htmx-todos

A teaching project that wraps **[jsonplaceholder.typicode.com/todos](https://jsonplaceholder.typicode.com/todos)**
as a server-rendered Go + HTMX app. Its structure mirrors
[`go-htmx-starter`](../go-htmx-starter): controllers → services → layout
rendering, but with an outbound HTTP API instead of Postgres. No database,
no auth, no npm.

```bash
make run              # templ generate + go run .  -> http://localhost:9001
PORT=9100 make run    # different port
JSONPLACEHOLDER_BASE_URL=http://localhost:9002 make run   # point at a fake upstream
```

Requires Go and the [templ CLI](https://templ.guide): `go install github.com/a-h/templ/cmd/templ@latest`

## What it wraps

| Our route | Upstream | Notes |
|---|---|---|
| `GET /` | — | 302 → `/todos` |
| `GET /todos` | `GET /todos` | 200 todos; `?q=` filters titles in memory |
| `GET /todos/{id}` | `GET /todos/{id}` | upstream 404 → our 404; other failures → 502 |
| `GET /public/*` | — | static files |

There is **no cache**: every request (including each debounced search
keystroke) calls upstream. That keeps the "wrapper" honest and makes the
service a clean place to add caching later (see exercises).

## What it demonstrates

| Concept | Where |
|---|---|
| One handler, two responses (full page vs fragment) | `pkg/util/render/layout.go` |
| DOM ids as template/code contracts | `internal/templates/constants/ids.go` |
| Layout shells chosen by `HX-Target` | `internal/templates/templatetargets/main.go` |
| Controller → service → outbound HTTP | `internal/controllers/public/todos.go` → `internal/services/todos.go` |
| Context propagation into outbound requests | `todos.go` (`http.NewRequestWithContext`) |
| Dependency injection (base URL + HTTP client) | `cmd/web.go`, `internal/services/service_manager.go` |
| Error mapping (404 vs 502) + HTMX toast | `internal/common/reqerror.go`, `todos.go` |
| Unit testing HTTP clients with `httptest` | `internal/services/todos_test.go` |

## Project structure

```
main.go                     dispatch -> cmd.Execute(name)
cmd/                        web() bootstrap: env -> client -> services -> AppContext -> router
internal/
  router/app_router.go      chi router: Recover middleware, mount public controllers
  middlewares/recover.go    panic -> 500, injects request-id logger into context
  context/                  AppContext (app-scoped DI) + ctx logger helpers
  services/
    todos.go                TodosService: HTTP wrapper + filtering + ErrNotFound
    service_manager.go      DI: holds base URL + HTTP client, lazily builds services
    todos_test.go           httptest-based tests (no network)
  controllers/public/       todos.go (list/detail), static.go (/public files)
  common/reqerror.go        HandleError: toast for HTMX, plain text otherwise
  templates/
    constants/ids.go        DOM ids shared by templates and Go code
    templatetargets/        "index" (full doc) and "page-layout" shells
    ui/layouts/             index, main (#main-content), base_page_layout
    ui/pages/               todos list page, todo detail page
    ui/partials/            nav, toast
    ui/tables/              todos table fragment
pkg/
  util/render/layout.go     THE switch: full page or fragment?
  model/todo.go             Todo struct matching upstream JSON
public/css/style.css        plain CSS
```

## The request lifecycle

### 1. First visit: `GET /todos` (no HTMX headers)

```
main.go -> cmd.Execute("web") -> web()             cmd/web.go
  -> base URL + http.Client{Timeout: 5s}
  -> ServiceManager, AppContext                    internal/services, internal/context
  -> router.BuildRouter: chi + Recover             internal/router/app_router.go
  -> GET / -> 302 /todos                           controllers/public/todos.go:redirect
  -> GET /todos -> todos.list
       service.List(ctx, "")                       internal/services/todos.go
         GET https://jsonplaceholder.typicode.com/todos   (ctx attached)
         decode 200 Todos, filter titles
       render.Layout(...)                          pkg/util/render/layout.go
         no HX-Request -> wrap in "index"          internal/templates/templatetargets/main.go
  <- COMPLETE HTML DOCUMENT (head, htmx.js, nav, page)
```

### 2. Nav / row click: `hx-get="/todos" hx-target="#main-content" hx-push-url`

Same handler, but `render.Layout` sees `HX-Target: main-content` (not in
the targets map) and falls through to render the **bare page component**:
just `<h1>Todos</h1>...`, swapped into `<main id="main-content">`.

### 3. Typing in the search box

```
keyup --(300ms delay)--> GET /todos?q=...    HX-Target: todos-table
todos.list fetches all todos upstream, filters by title
  -> content = tables.TodosTable(todos)       (not the whole page!)
  <- <table> only; swapped into #todos-table; input keeps focus
```

### 4. Detail: `GET /todos/1`

`todos.detail` calls the upstream detail endpoint (`/todos/1`) with the
request context, then renders `pages.TodoDetail`. Row clicks target
`#main-content`, so nav/header survive.

### 5. Errors

`common.HandleError` (`internal/common/reqerror.go`) adapts:
- HTMX request → status + `Hx-Toast: true` + rendered toast HTML →
  global listener in `layouts/index.templ` injects it into `#global-toast`
- normal request → plain `http.Error` text

Mapping: bad id or upstream 404 → **404**; network error, timeout, or
upstream 5xx → **502**.

## Verify

```bash
make test                                  # go test ./...
curl -s http://localhost:9001/todos        # full document
curl -s -H 'HX-Target: main-content' http://localhost:9001/todos   # fragment
curl -s -H 'HX-Target: todos-table' 'http://localhost:9001/todos?q=delectus'  # table only
```

## Where to go next

1. **Add caching** — wrap `TodosService.List` with a TTL cache (mutex +
   timestamp) so search doesn't hit upstream every keystroke. The service
   boundary is exactly the right seam.
2. **Swap the upstream for a real API** — add auth headers / API keys to
   the shared `http.Client` (e.g. a RoundTripper), or point
   `JSONPLACEHOLDER_BASE_URL` at your own service.
3. **Auth + per-route privileges** — see the starter's
   `internal/middlewares/auth.go` and `internal/builders/handler.go`.
4. **Database persistence** — see the starter's `internal/database/`
   (DAOs, go-jet, migrations) to replace the HTTP upstream.
# go-htmx-jobs
