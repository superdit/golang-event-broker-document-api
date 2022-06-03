package exception

func PanicIfNeeded(err interface{}) {
	if err != nil {
		panic(err)
	}
}

func PanicIfBadRequest(err interface{}) {
	if err != nil {
		panic(ValidationError{
			Status:    400,
			ErrorCode: 1000,
			Message:   "Invalid request",
		})
	}
}
