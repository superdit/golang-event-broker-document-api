package model

type ModelInsertRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ModelInsertResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ModelUpdateRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ModelUpdateResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ModelDeleteRequest struct {
	Id string `json:"id"`
}

type ModelEnityResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ModelListResponse struct {
	Total  int                  `json:"total"`
	Models []ModelEnityResponse `json:"models"`
}
