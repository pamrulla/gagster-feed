package user

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"

	"github.com/pamrulla/gagster-feed/database"
	"github.com/pamrulla/gagster-feed/models"
)

type UserRepo struct {
	Db  *gorm.DB
	usr models.UserHandlerInterface
}

func New() *UserRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.User{})
	uh := models.UserHandler{}
	return &UserRepo{Db: db, usr: uh}
}

func (ur *UserRepo) GetUsers(w http.ResponseWriter, r *http.Request) {
	var users models.Users
	err := ur.usr.GetUsers(ur.Db, &users)
	if err != nil {
		http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, users)
}

func (ur *UserRepo) Create(w http.ResponseWriter, r *http.Request) {
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
	err = ur.usr.CreateUser(ur.Db, &u)
	if err != nil {
		http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, "Successfully added new user")
}

func (ur *UserRepo) Update(w http.ResponseWriter, r *http.Request) {
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

	err = ur.usr.UpdateUser(ur.Db, &u)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}
	}
	render.JSON(w, r, u)
}

func (ur *UserRepo) Delete(w http.ResponseWriter, r *http.Request) {
	var u models.User
	vars := chi.URLParam(r, "user_id")

	err := ur.usr.DeleteUser(ur.Db, &u, vars)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}
	}
	render.JSON(w, r, "Successfully deleted user")
}

func (ur *UserRepo) Get(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "user_id")
	var u models.User
	err := ur.usr.GetUser(ur.Db, &u, vars)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}
	}
	render.JSON(w, r, u)
}

func (ur *UserRepo) Enable(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "user_id")
	var u models.User
	err := ur.usr.GetUser(ur.Db, &u, vars)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}
	}
	u.IsEnabled = true
	err = ur.usr.UpdateUser(ur.Db, &u)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}
	}
	render.JSON(w, r, "Successfully enabled user")
}

func (ur *UserRepo) Disable(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "user_id")
	var u models.User
	err := ur.usr.GetUser(ur.Db, &u, vars)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}
	}
	u.IsEnabled = false
	err = ur.usr.UpdateUser(ur.Db, &u)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}
	}
	render.JSON(w, r, "Successfully disabled user")
}
