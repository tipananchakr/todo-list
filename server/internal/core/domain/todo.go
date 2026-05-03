package domain

type Todo struct {
	ID        string `json:"id,omitempty"`
	UserID    string `json:"userId,omitempty"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}
