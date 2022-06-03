package model

type BackendInsertRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type BackendInsertResponse struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Url            string `json:"url"`
	TotalExchanges int    `json:"total_exchanges"`
	TotalQueues    int    `json:"total_queues"`
}

type BackendUpdateRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type BackendUpdateResponse struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Url            string `json:"url"`
	TotalExchanges int    `json:"total_exchanges"`
	TotalQueues    int    `json:"total_queues"`
}

type BackendDeleteRequest struct {
	Id string `json:"id"`
}

type BackendEnityResponse struct {
	Id             string `json:"id" bson:"_id"`
	Name           string `json:"name" bson:"name"`
	Url            string `json:"url"`
	TotalExchanges int    `json:"total_exchanges"`
	TotalQueues    int    `json:"total_queues"`
}

type BackendListResponse struct {
	Total    int                    `json:"total"`
	Backends []BackendEnityResponse `json:"backends"`
}
