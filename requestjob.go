package requestjob

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	file "github.com/Sotaneum/go-json-file"
)

func (job *RequestJob) toDate(targetTime time.Time) time.Time {
	return cronToDateWithoutSecond(job.Cron, targetTime)
}

// toDate : 실행할 날짜를 반환합니다.
func (job *RequestJob) toEqualDate(date time.Time) bool {
	return job.toDate(date).Equal(date)
}

// IsRun : Job를 실행해야하는 타임인지 여부를 반환합니다.
func (job *RequestJob) IsRun(t time.Time) bool {
	return !job.Inactive && job.toEqualDate(t)
}

// Save : Job를 파일로 저장합니다.
func (job *RequestJob) Save(path string) {
	f := file.File{Path: path, Name: job.ID + ".json"}
	f.Remove()
	f.SaveObject(job)
}

// Remove : Job를 삭제합니다.
func (job *RequestJob) Remove(path string) error {
	f := file.File{Path: path, Name: job.ID + ".json"}
	return f.Remove()
}

// BodyToCompile : body에 대한 추가적은 옵션이 필요하다면, 임베딩을 사용하세요.
func (job *RequestJob) BodyToCompile(body string) string {
	return body
}

// BodyToCompile : 기본 관리자를 수정하려면, 임베딩을 사용하세요.
func (job *RequestJob) GetDefaultOwner() string {
	return "admin"
}

func (job *RequestJob) bodyToBuffer(body string) (io.Reader, error) {
	var obj map[string]interface{}
	json.Unmarshal([]byte(body), &obj)

	jBody, err := json.Marshal(obj)

	if err != nil {
		return nil, ErrorJSONFormat
	}

	return bytes.NewBuffer(jBody), nil
}

// Run : Job를 실행합니다.
func (job *RequestJob) Run() interface{} {
	info := job.Reqeust
	body, bodyParseErr := job.bodyToBuffer(job.BodyToCompile(info.Body))

	if bodyParseErr != nil {
		return requestToObject(nil, bodyParseErr)
	}

	client := &http.Client{}

	methodStr := strings.ToUpper(info.Method)
	req, reqErr := http.NewRequest(methodStr, info.URL, body)

	if reqErr != nil {
		return requestToObject(nil, reqErr)
	}

	req.Header.Set("Content-Type", info.ContentType+"; charset=utf-8")
	res, resErr := client.Do(req)

	return requestToObject(res, resErr)
}

// SetOwner : Job의 주인을 설정합니다.
func (job *RequestJob) SetOwner(member string) {
	if member == "" {
		job.SetOwner(job.GetDefaultOwner())
		return
	}
	job.Admin.Owner = member
	for _, user := range job.Admin.Members {
		if user == member {
			return
		}
	}
	job.Admin.Members = append(job.Admin.Members, member)
}

// GetOwner : owner 정보를 반환합니다.
func (job *RequestJob) GetOwner() string {
	return job.Admin.Owner
}

// GetID : Job ID를 반환합니다.
func (job *RequestJob) GetID() string {
	return job.ID
}

// CreateID : ID를 생성합니다.
func (job *RequestJob) CreateID() error {
	hash := sha256.New()
	data, err := json.Marshal(job)

	if err != nil {
		return err
	}

	hash.Write(data)

	job.ID = hex.EncodeToString(hash.Sum(nil))

	return nil
}

// HasAuthorization : 주어진 멤버가 이 Job에 권한이 있는지 여부를 반환합니다.
func (job *RequestJob) HasAuthorization(member string) bool {
	if job.Admin.Owner == member {
		return true
	}

	members := job.Admin.Members

	for _, m := range members {
		if m == member {
			return true
		}
	}
	return false
}

// HasAdminAuthorization : Job를 관리자 수준까지 권한이 있는지 여부를 반한합니다.
func (job *RequestJob) HasAdminAuthorization(member string) bool {
	return job.Admin.Owner == member
}

// IsAvailability : 데이터가 유효성이 존재하는 지 여부를 반환합니다.
func (job *RequestJob) IsAvailability() bool {
	Owner := job.Admin.Owner
	members := job.Admin.Members

	for _, member := range members {
		if member == Owner {
			return true
		}
	}

	return false
}
