package db

import v1 "k8s.io/api/core/v1"

func GetResult(nodename map[string]v1.Node) []*Result {
	var result []*Result
	Dborm.QueryTable("result").OrderBy("-time").Limit(len(nodename)).All(&result)
	return result
}
