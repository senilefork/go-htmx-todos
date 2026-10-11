// Package context is this app's dependency container plus helpers for
// request-scoped values on context.Context (logger, ...). The name
// shadows the stdlib — the starter aliases it as `gctx`, we do the same.
package context

import "go-htmx-jobs/internal/services"

// AppContext is the application-scoped DI container: built once in
// cmd/web.go, passed to the router and embedded into every controller.
type AppContext struct {
	sm *services.ServiceManager
}

func NewAppContext(sm *services.ServiceManager) *AppContext {
	return &AppContext{sm: sm}
}

// SM lets controllers reach their services: jc.SM().TodosService()
func (a *AppContext) SM() *services.ServiceManager {
	return a.sm
}
