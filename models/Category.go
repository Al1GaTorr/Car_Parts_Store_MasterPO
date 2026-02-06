package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description" json:"description"`
	PartsList   []primitive.ObjectID `bson:"parts_list" json:"parts_list"`
}

func (cat *Category) AddCategory()                    {}
func (cat *Category) ListParts() []primitive.ObjectID { return cat.PartsList }
func (cat *Category) UpdateCategory()                 {}
