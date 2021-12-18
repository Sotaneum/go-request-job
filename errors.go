package requestjob

import "errors"

// ErrorNoAuthorization : 권한이 없는 경우 발생합니다.
var ErrorNoAuthorization = errors.New("권한이 없습니다")

// ErrorCantCreateJob : 데이터를 처리할 수 없을 때 (생성할 수 없을 떄) 발생합니다.
var ErrorCantCreateJob = errors.New("데이터를 처리할 수 없습니다")

var ErrorJSONFormat = errors.New("json format error")
