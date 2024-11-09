package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

type baseResponse struct {
	Message string `json:"message"`
}

type okResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewInternalServerError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	response := baseResponse{
		Message: message,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Fatal(err)
	}
}

func NewUnauthorized(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	response := baseResponse{
		Message: message,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Fatal(err)
	}
}

func NewUnprocessableEntity(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	response := baseResponse{
		Message: message,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Fatal(err)
	}
}

func NewBadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	response := baseResponse{
		Message: message,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Fatal(err)
	}
}

func NewNotFound(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNotFound)
	response := baseResponse{
		Message: message,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Fatal(err)
	}
}

func NewOKResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)
	response := okResponse{
		Message: message,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Fatal(err)
	}
}

func NewOKResponseWithData(w http.ResponseWriter, message string, data interface{}) {
	w.WriteHeader(http.StatusOK)
	response := okResponse{
		Message: message,
		Data:    data,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func NewOKResponseWithJson(w http.ResponseWriter, message string, data []byte) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func NewCreatedResponseWithData(w http.ResponseWriter, message string, data interface{}) {
	w.WriteHeader(http.StatusCreated)
	response := okResponse{
		Message: message,
		Data:    data,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func NewRedirect(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusSeeOther)
	response := baseResponse{
		Message: message,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Fatal(err)
	}
}
