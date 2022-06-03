package entity

type Model struct {
	Id   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Url  string `json:"url" bson:"url"`
}
