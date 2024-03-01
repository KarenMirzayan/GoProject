package main

import (
	"encoding/json"
	"github.com/KarenMirzayan/Project/pkg/messenger/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		return
	}
}

func (app *application) createUsersHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		First       string `json:"firstname"`
		Last        string `json:"lastname"`
		DateOfBirth string `json:"date_of_birth"`
		Login       string `json:"login"`
		Password    string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	users := &models.Users{
		First:       input.First,
		Last:        input.Last,
		DateOfBirth: input.DateOfBirth,
		Login:       input.Login,
		Password:    input.Password,
	}

	err = app.models.Users.Insert(users)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusCreated, users)
}

//func (app *application) getMenuHandler(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	param := vars["menuId"]
//
//	id, err := strconv.Atoi(param)
//	if err != nil || id < 1 {
//		app.respondWithError(w, http.StatusBadRequest, "Invalid menu ID")
//		return
//	}
//
//	menu, err := app.models.Menus.Get(id)
//	if err != nil {
//		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
//		return
//	}
//
//	app.respondWithJSON(w, http.StatusOK, menu)
//}

//func (app *application) updateUsersHandler(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	param := vars["usersId"]
//
//	id, err := strconv.Atoi(param)
//	if err != nil || id < 1 {
//		app.respondWithError(w, http.Stt atusBadRequest, "Invalid menu ID")
//		return
//	}
//
//	menu, err := app.models.Menus.Get(id)
//	if err != nil {
//		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
//		return
//	}
//
//	var input struct {
//		Title          *string `json:"title"`
//		Description    *string `json:"description"`
//		NutritionValue *uint   `json:"nutritionValue"`
//	}
//
//	err = app.readJSON(w, r, &input)
//	if err != nil {
//		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
//		return
//	}
//
//	if input.Title != nil {
//		menu.Title = *input.Title
//	}
//
//	if input.Description != nil {
//		menu.Description = *input.Description
//	}
//
//	if input.NutritionValue != nil {
//		menu.NutritionValue = *input.NutritionValue
//	}
//
//	err = app.models.Menus.Update(menu)
//	if err != nil {
//		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
//		return
//	}
//
//	app.respondWithJSON(w, http.StatusOK, menu)
//}

func (app *application) deleteUsersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["usersId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid menu ID")
		return
	}

	err = app.models.Users.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}
