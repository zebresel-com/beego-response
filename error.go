package beego_response

type Error struct {
	Code     int         `json:"code,omitempty"`
	Message  string      `json:"message"`
	UserInfo []*UserInfo `json:"userInfo,omitempty"`
}
