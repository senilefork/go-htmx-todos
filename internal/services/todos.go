package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go-htmx-jobs/pkg/model"
)

// DefaultBaseURL is used when JSONPLACEHOLDER_BASE_URL is not set.
const DefaultBaseURL = "https://jsonplaceholder.typicode.com"

// ErrNotFound is returned when the upstream API responds 404. Controllers
// map it to an HTTP 404; any other error becomes a 502.
var ErrNotFound = errors.New("todo not found")

// TodosService wraps the jsonplaceholder /todos endpoints. Controllers
// never make HTTP calls themselves — they go through this service, the
// same layering the starter uses for its DB-backed services.
type TodosService struct {
	baseURL string
	client  *http.Client
}

func NewTodosService(baseURL string, client *http.Client) *TodosService {
	return &TodosService{baseURL: baseURL, client: client}
}

// List fetches every todo, then filters by a case-insensitive substring
// match on the title. There is no title search upstream, so filtering
// happens here; with no cache this means one upstream call per search.
func (s *TodosService) List(ctx context.Context, query string) ([]model.Todo, error) {
	var todos []model.Todo
	if err := s.get(ctx, "/todos", &todos); err != nil {
		return nil, err
	}

	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return todos, nil
	}

	filtered := make([]model.Todo, 0, len(todos))
	for _, t := range todos {
		if strings.Contains(strings.ToLower(t.Title), q) {
			filtered = append(filtered, t)
		}
	}
	return filtered, nil
}

// GetByID fetches a single todo. Returns ErrNotFound on upstream 404.
func (s *TodosService) GetByID(ctx context.Context, id int) (model.Todo, error) {
	var todo model.Todo
	if err := s.get(ctx, fmt.Sprintf("/todos/%d", id), &todo); err != nil {
		return model.Todo{}, err
	}
	return todo, nil
}

// get performs a GET and decodes the JSON body. The context is attached
// to the outbound request, so if the browser disconnects (or a deadline
// fires) the upstream call is cancelled too.
func (s *TodosService) get(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("build request %s: %w", path, err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("call %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("call %s: unexpected status %d", path, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}
