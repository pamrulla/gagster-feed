package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/pamrulla/gagster-feed/models"
)

var users models.Users

func Init() {
	users = models.Users{
		models.User{Id: 1},
		models.User{Id: 2},
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, users)
}

func Create(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid data sent", http.StatusBadRequest)
		return
	}
	var u models.User
	err = json.Unmarshal(req, &u)
	if err != nil {
		http.Error(w, "Invalid data sent", http.StatusBadRequest)
		return
	}
	users = append(users, u)
	render.JSON(w, r, "Successfully added new user")
}

func Update(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid data sent", http.StatusBadRequest)
		return
	}
	var u models.User
	err = json.Unmarshal(req, &u)
	if err != nil {
		http.Error(w, "Invalid data sent", http.StatusBadRequest)
		return
	}

	isFound := false

	for i, a := range users {
		if a.Id == u.Id {
			users[i].First_Name = u.First_Name
			users[i].Last_Name = u.Last_Name
			isFound = true
			break
		}
	}
	if isFound {
		render.JSON(w, r, "Successfully updated user")
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "user_id")
	user_id, _ := strconv.Atoi(vars)

	isFound := false

	for i, a := range users {
		if a.Id == user_id {
			users = append(users[:i], users[i+1:]...)
			isFound = true
			break
		}
	}
	if isFound {
		render.JSON(w, r, "Successfully deleted user")
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "user_id")
	user_id, _ := strconv.Atoi(vars)

	for _, a := range users {
		if a.Id == user_id {
			render.JSON(w, r, a)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

func Enable(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "user_id")
	user_id, _ := strconv.Atoi(vars)

	for i, a := range users {
		if a.Id == user_id {
			users[i].IsEnabled = true
			render.JSON(w, r, "Successfully enabled user")
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

func Disable(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "user_id")
	user_id, _ := strconv.Atoi(vars)

	for i, a := range users {
		if a.Id == user_id {
			users[i].IsEnabled = false
			render.JSON(w, r, "Successfully disabled user")
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}
