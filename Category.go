package main

type Categories struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	PartsList   []string `json:"parts_list"`
}

func (cat *Categories) AddCategory()        {}
func (cat *Categories) ListParts() []string { return cat.PartsList }
func (cat *Categories) UpdateCategory()     {}
