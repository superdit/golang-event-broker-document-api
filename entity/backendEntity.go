package entity

type Backend struct {
	Id             string `json:"id" bson:"_id"`
	Name           string `json:"name" bson:"name"`
	Url            string `json:"url" bson:"url"`
	TotalExchanges int    `json:"total_exchanges" bson:"total_exchanges"`
	TotalQueues    int    `json:"total_queues" bson:"total_queues"`
}
