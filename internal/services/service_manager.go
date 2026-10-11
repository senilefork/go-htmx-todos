package services

import "net/http"

// ServiceManager holds one instance of every service (lazy singletons).
// It is built once at startup in cmd/web.go with the shared HTTP client
// and base URL, then handed to AppContext so controllers can reach it.
type ServiceManager struct {
	baseURL string
	client  *http.Client
	todos   *TodosService
}

func NewServiceManager(baseURL string, client *http.Client) *ServiceManager {
	return &ServiceManager{baseURL: baseURL, client: client}
}

func (m *ServiceManager) TodosService() *TodosService {
	if m.todos == nil {
		m.todos = NewTodosService(m.baseURL, m.client)
	}
	return m.todos
}

// INSERT SERVICE
