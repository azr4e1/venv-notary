package venv

type VenvNotRegisteredError struct {
	Message string
}

func (ae VenvNotRegisteredError) Error() string {
	return ae.Message
}

type MultipleVersionsError struct {
	Message string
}

func (mve MultipleVersionsError) Error() string {
	return mve.Message
}
