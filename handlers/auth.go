package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Nelwhix/numeris/pkg"
	"github.com/Nelwhix/numeris/pkg/requests"
	"github.com/Nelwhix/numeris/pkg/responses"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserResource struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes UserAttributes `json:"attributes"`
}

type UserAttributes struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Token    string `json:"token,omitempty"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	request, err := pkg.ParseRequestBody[requests.SignUp](r)
	if err != nil {
		responses.NewUnprocessableEntity(w, err.Error())
		return
	}

	err = h.Validator.Struct(request)
	if err != nil {
		responses.NewUnprocessableEntity(w, err.Error())
		return
	}

	_, err = h.Model.GetUserByEmail(r.Context(), request.Email)
	if err == nil {
		responses.NewBadRequest(w, "Email already taken")
		return
	}

	user, err := h.Model.InsertIntoUsers(r.Context(), request)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error creating user: %v", err.Error()))
		responses.NewInternalServerError(w, err.Error())
		return
	}

	response := UserResource{
		ID:   user.ID,
		Type: "user",
		Attributes: UserAttributes{
			Username: user.Username,
			Email:    user.Email,
		},
	}

	responses.NewCreatedResponseWithData(w, "User created successfully.", response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	request, err := pkg.ParseRequestBody[requests.Login](r)
	if err != nil {
		responses.NewUnprocessableEntity(w, err.Error())
		return
	}

	err = h.Validator.Struct(request)
	if err != nil {
		responses.NewUnprocessableEntity(w, err.Error())
		return
	}

	user, err := h.Model.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responses.NewBadRequest(w, "Email or Password is incorrect")
			return
		}

		responses.NewBadRequest(w, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		responses.NewBadRequest(w, "Email or Password is incorrect")
		return
	}

	token, err := pkg.GetOrCreateToken(h.Model, user.ID)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error creating token: %v", err.Error()))
		responses.NewInternalServerError(w, err.Error())
		return
	}

	response := UserResource{
		ID:   user.ID,
		Type: "user",
		Attributes: UserAttributes{
			Username: user.Username,
			Email:    user.Email,
			Token:    token,
		},
	}

	responses.NewOKResponseWithData(w, "Login success.", response)
}
