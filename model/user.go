// model/user.go
package model

// User represents a user in the system.
type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	// Password string `json:"password"`
}

// SignupRequest is the request body for the signup API.
type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest is the request body for the login API.
type LoginRequest struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignupResponse is the response body for the signup API.
type SignupResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// LoginResponse is the response body for the login API.
type LoginResponse struct {
	Token string `json:"token"`
}
