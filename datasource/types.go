package datasource

type Result struct {
	Id           string `db:"id"`
	Instance     string `db:"instance"`
	OriginIndex  string `db:"origin_index"`
	AutoIndex    string `db:"auto_index"`
	IforestIndex string `db:"iforest_index"`
	MergeIndex   string `db:"merge_index"`
	Time         string `db:"time"`
}
