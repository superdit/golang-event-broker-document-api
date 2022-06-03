package entity

type Event struct {
	Id                   string `json:"id" bson:"_id"`
	BackendPublisherId   string `json:"_backend_publisher_id" bson:"_backend_publisher_id"`
	BackendPublisherName string `json:"_backend_publisher_name" bson:"_backend_publisher_name"`
	ModelId              string `json:"_model_id" bson:"_model_id"`
	ModelName            string `json:"_model_name" bson:"_model_name"`
	ModelUrl             string `json:"_model_url" bson:"_model_url"`
	Exchange             string `json:"exchange" bson:"exchange"`
	TriggerAction        string `json:"trigger_action" bson:"trigger_action"`
	TotalSubscribers     int    `json:"total_subscribers" bson:"total_subscribers"`
	Endpoints            string `json:"endpoints" bson:"endpoints"`
	Description          string `json:"description" bson:"description"`
	SampleMessage        string `json:"sample_message" bson:"sample_message"`
}
