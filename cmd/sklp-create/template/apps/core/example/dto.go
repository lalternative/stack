package example

// CreateRequest is the POST /examples body.
type CreateRequest struct {
	Name string `json:"name"`
}

// RenameRequest is the PATCH /examples/:id body.
type RenameRequest struct {
	Name string `json:"name"`
}

// ExampleDTO is the read-model view returned to clients.
type ExampleDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
