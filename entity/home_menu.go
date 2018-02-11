package entity

type HomeMenu struct {
	ResId      string `json:"resId"`
	ResName    string `json:"resName"`
	ResUrl     string `json:"resUrl"`
	ResBgColor string `json:"resBgColor"`
	ResClass   string `json:"resClass"`
	ResImg     string `json:"resImg"`
	GroupId    string `json:"groupId"`
	ParentId   string `json:"parentId"`
	OpenTypeId string `json:"openTypeId"`
	NewIFrame  string `json:"newIFrame"`
	Method     string `json:"method"`
	ResAttr    string `json:"resAttr"`
}
