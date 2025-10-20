package internal

import (
	"encoding/json"

	pkgError "github.com/hgyowan/go-pkg-library/error"
)

func ParseMetadata(b []byte) (string, error) {
	var tmp map[string]interface{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return "", pkgError.Wrap(err)
	}

	return string(b), nil
}
