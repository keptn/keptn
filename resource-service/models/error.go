package models

// Error error
// swagger:model Error
type Error struct {

	// Error code
	Code int64 `json:"code,omitempty"`

	// Error message
	// Required: true
	Message string `json:"message"`
}

func (m *Error) Error() string {
	return m.Message
}
