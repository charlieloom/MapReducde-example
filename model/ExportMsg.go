package Model

import model "MapReduce/dal/model"

type ExportMsg struct {
	Productlist []*model.Product `json:"productlist"`
	File        string           `json:"file"`
	Row         int              `json:"row"`
}
