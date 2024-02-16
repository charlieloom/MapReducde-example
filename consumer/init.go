package consumer

import mapreduce "MapReduce/infra/mapReduce"

func Init() {
	m := mapreduce.GetMapReduce()
	go m.Reduce([]string{"query"}, Query)
	go m.Reduce([]string{"export"}, Export)
}
