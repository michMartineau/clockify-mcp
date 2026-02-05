package clockify

// User represents the current Clockify user.
type User struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	Name             string `json:"name"`
	DefaultWorkspace string `json:"defaultWorkspace"`
}