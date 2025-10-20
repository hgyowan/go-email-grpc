package internal

import (
	"encoding/json"
	pkgError "github.com/hgyowan/go-pkg-library/error"
)

func ParseMetadata(b []byte) (string, error) {
	var outer map[string]json.RawMessage
	err := json.Unmarshal(b, &outer)
	if err != nil {
		return "", pkgError.Wrap(err)
	}

	// 첫 번째 key의 value만 꺼냄
	for _, v := range outer {
		var inner map[string]interface{}
		err = json.Unmarshal(v, &inner)
		if err != nil {
			return "", pkgError.Wrap(err)
		}

		// 출력 확인
		resultBytes, _ := json.Marshal(inner)

		return string(resultBytes), nil
	}

	return "", pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.WrongParam)
}
