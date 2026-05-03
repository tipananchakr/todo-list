package domain

type Todo struct {
	ID        string `json:"id,omitempty"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}
