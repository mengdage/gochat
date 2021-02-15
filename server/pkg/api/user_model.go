package api

// User represents a newly created user.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
