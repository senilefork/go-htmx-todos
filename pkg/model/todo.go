package model

// Todo is the shape returned by jsonplaceholder.typicode.com/todos.
// Field names match the JSON so the service can decode directly.
type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
