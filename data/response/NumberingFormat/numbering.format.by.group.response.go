package numberinggroup

type NumberingFormatByGroupResponse struct {
	Group   string   `json:"group"`
	Formats []Format `json:"formats"`
}

type Format struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
