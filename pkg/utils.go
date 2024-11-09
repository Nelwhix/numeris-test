package pkg

import (
	"encoding/json"
	"fmt"
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

func FormatMoneyToUsd(amount *int64) string {
	if amount == nil {
		return "$0.00"
	}

	amountInUsd := float64(*amount) / 100.0

	return fmt.Sprintf("$%.2f", amountInUsd)
}
