package beego_response

type Paging struct {
	First        string `json:"first"`
	Last         string `json:"last"`
	Previous     string `json:"previous"`
	Next         string `json:"next"`
	RecordsTotal int    `json:"totalRecords"`
	RecordsPage  int    `json:"limit"`
	Pages        int    `json:"pages"`
	CurrentPage  int    `json:"currentPage"`
}
