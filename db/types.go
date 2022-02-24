package db

import (
	"time"
)

type Result struct {
	Id           string    `orm:"column(id);pk"`
	Instance     string    `orm:"column(instance)"`
	OriginIndex  string    `orm:"column(origin_index)"`
	AutoIndex    string    `orm:"column(auto_index)"`
	IforestIndex string    `orm:"column(iforest_index)"`
	MergeIndex   string    `orm:"column(merge_index)"`
	Time         time.Time `orm:"column(time)"`
}
