package requestjob

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	gocreatefolder "github.com/Sotaneum/go-create-folder"
	file "github.com/Sotaneum/go-json-file"
	ktime "github.com/Sotaneum/go-kst-time"
)

// ErrorNoAuthorization : 권한이 없는 경우 발생합니다.
var ErrorNoAuthorization = errors.New("권한이 없습니다")

// ErrorCantCreateJob : 데이터를 처리할 수 없을 때 (생성할 수 없을 떄) 발생합니다.
var ErrorCantCreateJob = errors.New("데이터를 처리할 수 없습니다")

// NewByFile : File로 부터 Job 객체를 생성합니다.
func NewByFile(path, name, owner string) (*RequestJob, error) {
	f := file.File{Path: path, Name: name}

	data := f.Load()

	if data == "" {
		return nil, ErrorCantCreateJob
	}

	return NewByJSON(data, owner)
}

// NewByJSON : json 데이터로 부터 Job 객체를 생성합니다.
func NewByJSON(data, owner string) (*RequestJob, error) {
	job := RequestJob{}

	json.Unmarshal([]byte(data), &job)

	if job.Admin.Owner == "" {
		job.SetOwner(owner)
	}

	if !job.IsAvailability() {
		return nil, ErrorNoAuthorization
	}

	if job.CreateDate == "" {
		job.CreateDate = ktime.GetNowWithSecond().String()
	}

	if job.ID == "" {
		if err := job.CreateID(); err != nil {
			return nil, ErrorCantCreateJob
		}
	}

	return &job, nil
}

// NewList : 폴더에 있는 데이터를 모두 Job객체로 만들어 반환합니다.
func NewList(path string) ([]*RequestJob, error) {
	jobList := []*RequestJob{}

	createFolderErr := gocreatefolder.CreateFolder(path, 0755)

	if createFolderErr != nil {
		return jobList, createFolderErr
	}

	files, err := ioutil.ReadDir(path)

	if err != nil {
		return jobList, err
	}

	for _, f := range files {
		job, jobErr := NewByFile(path, f.Name(), "")

		if jobErr == nil {
			jobList = append(jobList, job)
		}
	}

	return jobList, nil
}

func New() *RequestJob {
	job := RequestJob{}
	job.CreateID()
	return &job
}
