package model

type EventInsertRequest struct {
	BackendPublisherId string `json:"_backend_publisher_id"`
	ModelId            string `json:"_model_id"`
	Exchange           string `json:"exchange"`
	TriggerAction      string `json:"trigger_action"`
	Endpoints          string `json:"endpoints"`
	Description        string `json:"description"`
	SampleMessage      string `json:"sample_message"`
}

type EventInsertResponse struct {
	Id                   string `json:"_id"`
	BackendPublisherName string `json:"_backend_publisher_name"`
	ModelName            string `json:"_model_name"`
	EventInsertRequest
}

type EventUpdateRequest struct {
	EventInsertResponse
}

type EventUpdateResponse struct {
	EventInsertResponse
}

type EventDeleteRequest struct {
	Id string `json:"id"`
}

type EventEnityResponse struct {
	Id                   string `json:"id"`
	BackendPublisherId   string `json:"_backend_publisher_id"`
	BackendPublisherName string `json:"_backend_publisher_name"`
	ModelId              string `json:"_model_id"`
	ModelName            string `json:"_model_name"`
	ModelUrl             string `json:"_model_url"`
	Exchange             string `json:"exchange"`
	TriggerAction        string `json:"trigger_action"`
	TotalSubscribers     int    `json:"total_subscribers"`
	Endpoints            string `json:"endpoints"`
	Description          string `json:"description"`
	SampleMessage        string `json:"sample_message"`
}

type EventListResponse struct {
	Total  int                  `json:"total"`
	Events []EventEnityResponse `json:"events"`
}
