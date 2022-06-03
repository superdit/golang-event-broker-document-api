package repository

import "event-broker-document-api/entity"

type ModelRepository interface {
	List(query string) (ents []entity.Model)
	FindByName(name string) (ents entity.Model)
	FindById(id string) (ent entity.Model)
	FindNameExceptId(name string, id string) (ent entity.Model)
	Insert(ent entity.Model)
	Update(ent entity.Model)
	Delete(entId string)
}
