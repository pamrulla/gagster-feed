package gag

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/pamrulla/gagster-feed/models"
)

var gags models.Gags

func Init() {
	gags = models.Gags{
		models.Gag{Id: 1, User_Id: 1, Price: 100, Path: "path1"},
		models.Gag{Id: 2, User_Id: 2, Price: 110, Path: "path2"},
	}
}

func getUserIdFromRequest(r *http.Request) int {
	vars := chi.URLParam(r, "user_id")
	user_id, _ := strconv.Atoi(vars)
	return user_id
}

func GetGags(w http.ResponseWriter, r *http.Request) {
	user_id := getUserIdFromRequest(r)
	var res models.Gags
	for _, g := range gags {
		if g.User_Id == user_id {
			res = append(res, g)
		}
	}
	render.JSON(w, r, res)
}

func Create(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid data sent", http.StatusBadRequest)
		return
	}
	var u models.Gag
	err = json.Unmarshal(req, &u)
	if err != nil {
		http.Error(w, "Invalid data sent", http.StatusBadRequest)
		return
	}

	user_id := getUserIdFromRequest(r)
	u.User_Id = user_id

	gags = append(gags, u)
	render.JSON(w, r, "Successfully added new gag")
}

func Update(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid data sent", http.StatusBadRequest)
		return
	}
	var u models.Gag
	err = json.Unmarshal(req, &u)
	if err != nil {
		http.Error(w, "Invalid data sent", http.StatusBadRequest)
		return
	}

	isFound := false
	user_id := getUserIdFromRequest(r)

	for i, a := range gags {
		if a.Id == u.Id && a.User_Id == user_id {
			gags[i].Path = u.Path
			gags[i].Price = u.Price
			isFound = true
			break
		}
	}
	if isFound {
		render.JSON(w, r, "Successfully updated gag")
	} else {
		http.Error(w, "Gag not found", http.StatusNotFound)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "gag_id")
	gag_id, _ := strconv.Atoi(vars)

	user_id := getUserIdFromRequest(r)
	isFound := false

	for i, a := range gags {
		if a.Id == gag_id && user_id == a.User_Id {
			gags = append(gags[:i], gags[i+1:]...)
			isFound = true
			break
		}
	}
	if isFound {
		render.JSON(w, r, "Successfully deleted gag")
	} else {
		http.Error(w, "Gag not found", http.StatusNotFound)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "gag_id")
	gag_id, _ := strconv.Atoi(vars)
	user_id := getUserIdFromRequest(r)

	for _, a := range gags {
		if a.Id == gag_id && a.User_Id == user_id {
			render.JSON(w, r, a)
			return
		}
	}
	http.Error(w, "Gag not found", http.StatusNotFound)
}

func Enable(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "gag_id")
	gag_id, _ := strconv.Atoi(vars)

	for i, a := range gags {
		if a.Id == gag_id {
			gags[i].IsEnabled = true
			render.JSON(w, r, "Successfully enabled gag")
			return
		}
	}
	http.Error(w, "Gag not found", http.StatusNotFound)
}

func Disable(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "gag_id")
	gag_id, _ := strconv.Atoi(vars)

	for i, a := range gags {
		if a.Id == gag_id {
			gags[i].IsEnabled = false
			render.JSON(w, r, "Successfully disabled gag")
			return
		}
	}
	http.Error(w, "Gag not found", http.StatusNotFound)
}
