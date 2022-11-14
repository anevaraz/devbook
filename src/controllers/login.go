package controllers

import (
	"devbook/src/auth"
	"devbook/src/database"
	"devbook/src/models"
	"devbook/src/repositories"
	"devbook/src/response"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var loginUser models.User
	if err = json.Unmarshal(body, &loginUser); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connection()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.Users(db)
	user, err := repository.FindOneByEmail(loginUser.Email)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err = auth.Verify(user.Password, loginUser.Password); err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Write([]byte(token))
}
