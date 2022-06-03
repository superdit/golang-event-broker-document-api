package entity

type Subscriber struct {
	Id                   string `json:"id" bson:"_id"`
	EventId              string `json:"_event_id" bson:"_event_id"`
	Exchange             string `json:"_exchange" bson:"_exchange"`
	ExchangeModelUrl     string `json:"_exchange_model_url" bson:"_exchange_model_url"`
	BackendSubcriberId   string `json:"_backend_subscriber_id" bson:"_backend_subscriber_id"`
	BackendSubcriberName string `json:"_backend_subscriber_name" bson:"_backend_subscriber_name"`
	Queue                string `json:"queue" bson:"queue"`
	PostAction           string `json:"post_action" bson:"post_action"`
}
