package model

type EmptyResponse struct {
}

type WebResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type WebResponseError struct {
	Error ErrorResponse `json:"error"`
}
