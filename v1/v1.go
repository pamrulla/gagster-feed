package v1

import (
	"net/http"

	"github.com/pamrulla/gagster-feed/v1/gag"
	"github.com/pamrulla/gagster-feed/v1/heart"
	"github.com/pamrulla/gagster-feed/v1/user"
)

var ur *user.UserRepo
var gr *gag.GagRepo
var hr *heart.HeartRepo

func Init() {
	ur = user.New()
	gr = gag.New()
	hr = heart.New()
}

// User Table Hanlders
func GetUsers(w http.ResponseWriter, r *http.Request) {
	ur.GetUsers(w, r)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ur.Create(w, r)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	ur.Update(w, r)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ur.Delete(w, r)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	ur.Get(w, r)
}

func EnableUser(w http.ResponseWriter, r *http.Request) {
	ur.Enable(w, r)
}

func DisableUser(w http.ResponseWriter, r *http.Request) {
	ur.Disable(w, r)
}

// Gag Table Handler
func GetGags(w http.ResponseWriter, r *http.Request) {
	gr.GetGags(w, r)
}

func CreateGag(w http.ResponseWriter, r *http.Request) {
	gr.Create(w, r)
}

func UpdateGag(w http.ResponseWriter, r *http.Request) {
	gr.Update(w, r)
}

func DeleteGag(w http.ResponseWriter, r *http.Request) {
	gr.Delete(w, r)
}

func GetGag(w http.ResponseWriter, r *http.Request) {
	gr.Get(w, r)
}

func EnableGag(w http.ResponseWriter, r *http.Request) {
	gr.Enable(w, r)
}

func DisableGag(w http.ResponseWriter, r *http.Request) {
	gr.Disable(w, r)
}

// Hearts Table Handler
func GetHearts(w http.ResponseWriter, r *http.Request) {
	hr.GetHearts(w, r)
}

func CreateHeart(w http.ResponseWriter, r *http.Request) {
	hr.Create(w, r)
}

func DeleteHeart(w http.ResponseWriter, r *http.Request) {
	hr.Delete(w, r)
}
