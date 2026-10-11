package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newUpstream is a fake jsonplaceholder: /todos returns a small list,
// /todos/1 a single todo, and any other /todos/{id} is a 404.
func newUpstream(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{"userId":1,"id":1,"title":"write the plan","completed":false},
			{"userId":1,"id":2,"title":"ship the feature","completed":true},
			{"userId":2,"id":3,"title":"write docs","completed":false}
		]`)
	})
	mux.HandleFunc("/todos/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"userId":1,"id":1,"title":"write the plan","completed":false}`)
	})
	return httptest.NewServer(mux)
}

func TestList(t *testing.T) {
	srv := newUpstream(t)
	defer srv.Close()

	svc := NewTodosService(srv.URL, srv.Client())
	got, err := svc.List(context.Background(), "")
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("want 3 todos, got %d", len(got))
	}
	if got[0].Title != "write the plan" {
		t.Errorf("want first title %q, got %q", "write the plan", got[0].Title)
	}
}

func TestListFiltersByTitle(t *testing.T) {
	srv := newUpstream(t)
	defer srv.Close()

	svc := NewTodosService(srv.URL, srv.Client())
	got, err := svc.List(context.Background(), "DOCS") // case-insensitive
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(got) != 1 || got[0].ID != 3 {
		t.Fatalf("want only todo 3, got %+v", got)
	}
}

func TestGetByID(t *testing.T) {
	srv := newUpstream(t)
	defer srv.Close()

	svc := NewTodosService(srv.URL, srv.Client())
	got, err := svc.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}
	if got.ID != 1 || got.UserID != 1 {
		t.Errorf("unexpected todo: %+v", got)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	srv := newUpstream(t)
	defer srv.Close()

	svc := NewTodosService(srv.URL, srv.Client())
	_, err := svc.GetByID(context.Background(), 999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestListUpstreamError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	svc := NewTodosService(srv.URL, srv.Client())
	_, err := svc.List(context.Background(), "")
	if err == nil {
		t.Fatal("want an error for upstream 500, got nil")
	}
}
