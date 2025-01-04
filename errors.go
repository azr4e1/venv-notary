package venv

type VenvNotRegisteredError struct {
	Message string
}

func (ae VenvNotRegisteredError) Error() string {
	return ae.Message
}
