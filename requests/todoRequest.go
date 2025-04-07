package requests

type TodoRequest struct {
	Title  string `json:"title" validate:"required,min=2,max=10"`
	IsDone bool   `json:"is_done"`
}
