package pkg

import (
	"encoding/json"
	"io"
	"net/http"
)

func ParseRequestBody[T any](r *http.Request) (T, error) {
	var requestData T

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return requestData, err
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		return requestData, err
	}

	return requestData, nil
}
