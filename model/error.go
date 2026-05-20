package model

// ─── Domain Errors ─────────────────────────────────────────

type NotFoundError struct {
	ID int64
}

func (e *NotFoundError) Error() string { return "not found" }

type ValidationError struct {
	Code string
}

func (e *ValidationError) Error() string { return e.Code }
