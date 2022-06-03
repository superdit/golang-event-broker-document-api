package exception

type ValidationError struct {
	Message   string
	Status    int
	ErrorCode int
}

func (validationError ValidationError) Error() string {
	return validationError.Message
}
