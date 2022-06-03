package model

type SubscriberInsertRequest struct {
	EventId            string `json:"_event_id"`
	BackendSubcriberId string `json:"_backend_subscriber_id"`
	Queue              string `json:"queue"`
	PostAction         string `json:"post_action"`
}

type SubscriberInsertResponse struct {
	Id                   string `json:"_id"`
	EventId              string `json:"_event_id"`
	Exchange             string `json:"_exchange"`
	ExchangeModelUrl     string `json:"_exchange_model_url"`
	BackendSubcriberId   string `json:"_backend_subscriber_id"`
	BackendSubcriberName string `json:"_backend_subscriber_name"`
	Queue                string `json:"queue"`
	PostAction           string `json:"post_action"`
}

type SubscriberUpdateRequest struct {
	Id string `json:"_id"`
	SubscriberInsertRequest
}

type SubscriberUpdateResponse struct {
	SubscriberInsertResponse
}

type SubscriberDeleteRequest struct {
	Id string `json:"_id"`
}

type SubscriberEnityResponse struct {
	Id                   string `json:"_id"`
	EventId              string `json:"_event_id"`
	Exchange             string `json:"_exchange"`
	ExchangeModelUrl     string `json:"_exchange_model_url"`
	BackendSubcriberId   string `json:"_backend_subscriber_id"`
	BackendSubcriberName string `json:"_backend_subscriber_name"`
	Queue                string `json:"queue"`
	PostAction           string `json:"post_action"`
}

type SubscriberListResponse struct {
	Total       int                       `json:"total"`
	Subscribers []SubscriberEnityResponse `json:"subscribers"`
}
