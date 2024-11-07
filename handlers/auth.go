package handlers

import (
	"github.com/Nelwhix/numeris/pkg"
	"github.com/Nelwhix/numeris/pkg/requests"
	"github.com/Nelwhix/numeris/pkg/responses"
	"net/http"
)

type User struct {
	ID string `json:"id"`
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
		responses.NewInternalServerError(w, err.Error())
		return
	}

	response := User{
		ID:   user.ID,
		Type: "user",
		Attributes: responses.UserAttributes{
			Username: user.Username,
			Email:    user.Email,
		},
	}

	responses.NewCreatedResponseWithData(w, "User created successfully.", response)
}
