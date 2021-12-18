package requestjob

// ResponseData : request 이후 반환되는 데이터 포맷입니다.
type ResponseData struct {
	Code   int
	Status string
	Body   string
	Header string
}

type request struct {
	URL         string `json:"url"`
	Method      string `json:"method"`
	ContentType string `json:"contentType"`
	Body        string `json:"body"`
	Header      string `json:"header"`
}

type admin struct {
	Owner   string   `json:"owner"`
	Members []string `json:"members"`
}

type extra struct {
	Type string `json:"type"`
}

type RequestJob struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Cron       string  `json:"cron"`
	Reqeust    request `json:"request"`
	Inactive   bool    `json:"inactive"`
	CreateDate string  `json:"createDate"`
	Admin      admin   `json:"admin"`
	Extra      extra   `json:"extra"`
}
