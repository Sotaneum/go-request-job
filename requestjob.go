package requestjob

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	ktime "github.com/Sotaneum/go-kst-time"
	runnerjob "github.com/Sotaneum/go-runner-job"
	"github.com/gorhill/cronexpr"
)

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

type RequestJob struct {
	runnerjob.BaseJob
	Name    string  `json:"name"`
	Cron    string  `json:"cron"`
	Reqeust request `json:"request"`
}

var ErrorJSONFormat = errors.New("json format error")

func (job *RequestJob) requestToObject(res *http.Response, err error) *ResponseData {
	data := ResponseData{}

	if res == nil {
		data.Code = 500
		data.Status = "500 Request Error"
		data.Body = err.Error()
		return &data
	}

	header, headerErr := json.Marshal(res.Header)

	if headerErr != nil {
		data.Header = err.Error()
	} else {
		data.Header = string(header)
	}

	data.Code = res.StatusCode
	data.Status = res.Status

	body, er := ioutil.ReadAll(res.Body)

	if er != nil {
		data.Body = er.Error()
	} else {
		data.Body = string(body)
	}

	return &data
}

func (job *RequestJob) toDate(targetTime time.Time) time.Time {
	return ktime.ToDateByTime(cronexpr.MustParse(job.Cron).Next(targetTime.Add(-1)))
}

// toDate : 실행할 날짜를 반환합니다.
func (job *RequestJob) toEqualDate(date time.Time) bool {
	return job.toDate(date).Equal(date)
}

// IsRun : Job를 실행해야하는 타임인지 여부를 반환합니다.
func (job *RequestJob) IsRun(t time.Time) bool {
	return job.Active && job.toEqualDate(t)
}

// BodyToCompile : body에 대한 추가적은 옵션이 필요하다면, 임베딩을 사용하세요.
func (job *RequestJob) BodyToCompile(body string) string {
	return body
}

// json일경우에만 JSON 만들기
func (job *RequestJob) bodyToBuffer() (io.Reader, error) {
	info := job.Reqeust
	body := job.BodyToCompile(info.Body)

	var obj map[string]interface{}
	json.Unmarshal([]byte(body), &obj)

	jBody, err := json.Marshal(obj)

	if err != nil {
		return nil, ErrorJSONFormat
	}

	if jBody == nil && body != "" {
		return nil, ErrorJSONFormat
	}

	return bytes.NewBuffer(jBody), nil
}

// Run : Job를 실행합니다.
func (job *RequestJob) Run() interface{} {
	info := job.Reqeust
	body, bodyParseErr := job.bodyToBuffer()

	if bodyParseErr != nil {
		return job.requestToObject(nil, bodyParseErr)
	}

	client := &http.Client{}

	methodStr := strings.ToUpper(info.Method)
	req, reqErr := http.NewRequest(methodStr, info.URL, body)

	if reqErr != nil {
		return job.requestToObject(nil, reqErr)
	}

	req.Header.Set("Content-Type", info.ContentType+"; charset=utf-8")
	res, resErr := client.Do(req)

	return job.requestToObject(res, resErr)
}
