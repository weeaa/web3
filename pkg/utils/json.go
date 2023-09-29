package utils

import (
	"encoding/json"
	"io"
)

func UnmarshalJSONToStruct[T any](respBody io.ReadCloser) (T, error) {
	var t T
	body, err := io.ReadAll(respBody)
	if err != nil {
		return t, err
	}
	if err = json.Unmarshal(body, &t); err != nil {
		return t, err
	}
	return t, nil
}
