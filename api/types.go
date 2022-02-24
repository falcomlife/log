package api

type Chart struct {
	XAxis  XAxis    `json:"xAxis"`
	Series []Series `json:"series"`
}
type XAxis struct {
	Data []string `jaon:"data"`
}
type Series struct {
	Data     []string `json:"data"`
	MarkArea MarkArea `json:"markArea"`
}
type MarkArea struct {
	Data [][]Data `json:"data"`
}
type Data struct {
	Name  string `json:"name"`
	XAxis string `json:xAxis`
}
