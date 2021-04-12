package gag

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/pamrulla/gagster-feed/database"
	"github.com/pamrulla/gagster-feed/models"
	"gorm.io/gorm"
)

type GagRepo struct {
	Db *gorm.DB
}

func New() *GagRepo {
	db := database.InitDb()
	db.AutoMigrate(models.Gag{})
	return &GagRepo{Db: db}
}

func (gr *GagRepo) checkErr(err error, w http.ResponseWriter) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Gag not found", http.StatusNotFound)
	} else {
		http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
	}
}

func (gr *GagRepo) GetGags(w http.ResponseWriter, r *http.Request) {
	user_id := chi.URLParam(r, "user_id")
	var res models.Gags
	err := models.GetGags(gr.Db, &res, user_id)
	if err != nil {
		gr.checkErr(err, w)
		return
	}
	render.JSON(w, r, res)
}

func (gr *GagRepo) Create(w http.ResponseWriter, r *http.Request) {
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

	user_id := chi.URLParam(r, "user_id")
	u.User_Id, _ = strconv.Atoi(user_id)
	err = models.CreateGag(gr.Db, &u)
	if err != nil {
		gr.checkErr(err, w)
		return
	}
	render.JSON(w, r, "Successfully added new gag")
}

func (gr *GagRepo) Update(w http.ResponseWriter, r *http.Request) {
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

	err = models.UpdateGag(gr.Db, &u)
	if err != nil {
		gr.checkErr(err, w)
		return
	}
	render.JSON(w, r, u)
}

func (gr *GagRepo) Delete(w http.ResponseWriter, r *http.Request) {
	gag_id := chi.URLParam(r, "gag_id")

	var gag models.Gag
	err := models.DeleteGag(gr.Db, &gag, gag_id)
	if err != nil {
		gr.checkErr(err, w)
		return
	}

	render.JSON(w, r, "Successfully deleted gag")
}

func (gr *GagRepo) Get(w http.ResponseWriter, r *http.Request) {
	gag_id := chi.URLParam(r, "gag_id")
	var g models.Gag
	err := models.GetGag(gr.Db, &g, gag_id)
	if err != nil {
		gr.checkErr(err, w)
		return
	}
	render.JSON(w, r, g)
}

func (gr *GagRepo) Enable(w http.ResponseWriter, r *http.Request) {
	gag_id := chi.URLParam(r, "gag_id")

	var g models.Gag
	err := models.GetGag(gr.Db, &g, gag_id)
	if err != nil {
		gr.checkErr(err, w)
		return
	}
	g.IsEnabled = true
	err = models.UpdateGag(gr.Db, &g)
	if err != nil {
		gr.checkErr(err, w)
		return
	}
	render.JSON(w, r, "Successfully enabled gag")
}

func (gr *GagRepo) Disable(w http.ResponseWriter, r *http.Request) {
	gag_id := chi.URLParam(r, "gag_id")

	var g models.Gag
	err := models.GetGag(gr.Db, &g, gag_id)
	if err != nil {
		gr.checkErr(err, w)
		return
	}
	g.IsEnabled = false
	err = models.UpdateGag(gr.Db, &g)
	if err != nil {
		gr.checkErr(err, w)
		return
	}
	render.JSON(w, r, "Successfully disabled gag")
}
