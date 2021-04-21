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
	"github.com/pamrulla/gagster-feed/helpers"
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
		helpers.NewError(w, r, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, users)
}

func (ur *UserRepo) Create(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.NewError(w, r, "Invalid data sent", http.StatusBadRequest)
		return
	}
	var u models.User
	err = json.Unmarshal(req, &u)
	if err != nil {
		helpers.NewError(w, r, "Invalid data sent", http.StatusBadRequest)
		return
	}

	isUserExists := ur.usr.IsUserExists(ur.Db, u.Email)
	if isUserExists {
		helpers.NewError(w, r, "User with same email id already exists", http.StatusConflict)
		return
	}
	err = ur.usr.CreateUser(ur.Db, &u)
	if err != nil {
		helpers.NewError(w, r, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, "Successfully added new user")
}

func (ur *UserRepo) Update(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.NewError(w, r, "Invalid data sent", http.StatusBadRequest)
		return
	}
	var u models.User
	err = json.Unmarshal(req, &u)
	if err != nil {
		helpers.NewError(w, r, "Invalid data sent", http.StatusBadRequest)
		return
	}

	err = ur.usr.UpdateUser(ur.Db, &u)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.NewError(w, r, "User not found", http.StatusNotFound)
			return
		} else {
			helpers.NewError(w, r, "Something went wrong, please try again", http.StatusInternalServerError)
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
			helpers.NewError(w, r, "User not found", http.StatusNotFound)
			return
		} else {
			helpers.NewError(w, r, "Something went wrong, please try again", http.StatusInternalServerError)
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
			helpers.NewError(w, r, "User not found", http.StatusNotFound)
			return
		} else {
			helpers.NewError(w, r, "Something went wrong, please try again", http.StatusInternalServerError)
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
			helpers.NewError(w, r, "User not found", http.StatusNotFound)
			return
		} else {
			helpers.NewError(w, r, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}
	}
	u.IsEnabled = true
	err = ur.usr.UpdateUser(ur.Db, &u)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.NewError(w, r, "User not found", http.StatusNotFound)
			return
		} else {
			helpers.NewError(w, r, "Something went wrong, please try again", http.StatusInternalServerError)
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
			helpers.NewError(w, r, "User not found", http.StatusNotFound)
			return
		} else {
			helpers.NewError(w, r, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}
	}
	u.IsEnabled = false
	err = ur.usr.UpdateUser(ur.Db, &u)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.NewError(w, r, "User not found", http.StatusNotFound)
			return
		} else {
			helpers.NewError(w, r, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}
	}
	render.JSON(w, r, "Successfully disabled user")
}

func (ur *UserRepo) LogIn(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.NewError(w, r, "Invalid data sent", http.StatusBadRequest)
		return
	}
	var cred map[string]string

	err = json.Unmarshal(req, &cred)
	if err != nil {
		helpers.NewError(w, r, "Invalid data sent", http.StatusBadRequest)
		return
	}
	_, e_ok := cred["email"]
	_, p_ok := cred["password"]
	if !e_ok || !p_ok {
		helpers.NewError(w, r, "Invalid data sent", http.StatusBadRequest)
		return
	}
	var u models.User
	err = ur.usr.LoginUser(ur.Db, &u, cred["email"], cred["password"])
	if err != nil {
		helpers.NewError(w, r, "Failed to authenticate, please check your email and password", http.StatusUnauthorized)
		return
	}
	result := make(map[string]interface{})
	result["user"] = u
	result["auth_token"] = helpers.GenerateTokenString(u.Email, int(u.ID))
	render.JSON(w, r, result)
}
