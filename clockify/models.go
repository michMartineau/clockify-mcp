package clockify

// User represents the current Clockify user.
type User struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	Name             string `json:"name"`
	DefaultWorkspace string `json:"defaultWorkspace"`
}

// Project represents a Clockify project.
type Project struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ClientID   string `json:"clientId"`
	ClientName string `json:"clientName"`
	Archived   bool   `json:"archived"`
	Billable   bool   `json:"billable"`
	Color      string `json:"color"`
}

// Tag represents a Clockify tag.
type Tag struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	WorkspaceID string `json:"workspaceId"`
	Archived    bool   `json:"archived"`
}

// TimeEntry represents a Clockify time entry.
type TimeEntry struct {
	ID           string       `json:"id"`
	Description  string       `json:"description"`
	ProjectID    string       `json:"projectId"`
	TagIDs       []string     `json:"tagIds"`
	Billable     bool         `json:"billable"`
	TimeInterval TimeInterval `json:"timeInterval"`
}

// TimeInterval represents the start/end times of a time entry.
type TimeInterval struct {
	Start    string `json:"start"`    // ISO 8601
	End      string `json:"end"`      // ISO 8601, empty if running
	Duration string `json:"duration"` // ISO 8601 duration
}

// StartTimerRequest represents the request body for starting a timer.
type StartTimerRequest struct {
	Start       string   `json:"start"`
	Description string   `json:"description,omitempty"`
	ProjectID   string   `json:"projectId,omitempty"`
	TagIDs      []string `json:"tagIds,omitempty"`
	Billable    bool     `json:"billable,omitempty"`
}