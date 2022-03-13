package vo

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Pagination struct {
	CurPage  int         `json:"curPage"`
	TotalNum int64       `json:"totalNum"`
	Items    interface{} `json:"items"`
}
