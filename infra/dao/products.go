package dao

import (
	model "MapReduce/dal/model"
	"MapReduce/dal/query"
	model2 "MapReduce/model"
	"context"
)

func GetAllproducts(condition *model2.Condition, offest int, limit int) ([]*model.Product, error) {
	p := query.Product
	query := p.WithContext(context.Background()).Where()
	if condition.Id != 0 {
		query = query.Where(p.ID.Eq(int32(condition.Id)))
	}
	if condition.Name != "" {
		query = query.Where(p.Name.Eq(condition.Name))
	}
	if condition.Category != "" {
		query = query.Where(p.Category.Eq(string(condition.Category)))
	}
	products, err := query.Limit(limit).Offset(offest).Find()
	if err != nil {
		return nil, err
	}
	return products, nil
}
