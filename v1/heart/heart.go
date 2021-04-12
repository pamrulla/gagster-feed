package heart

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/pamrulla/gagster-feed/database"
	"github.com/pamrulla/gagster-feed/models"
	"gorm.io/gorm"
)

type HeartRepo struct {
	Db *gorm.DB
}

func New() *HeartRepo {
	db := database.InitDb()
	db.AutoMigrate(models.Heart{})
	return &HeartRepo{Db: db}
}

func (gr *HeartRepo) checkErr(err error, w http.ResponseWriter) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Heart not found", http.StatusNotFound)
	} else {
		http.Error(w, "Something went wrong, please try again", http.StatusInternalServerError)
	}
}

func (gr *HeartRepo) GetHearts(w http.ResponseWriter, r *http.Request) {
	gag_id := chi.URLParam(r, "gag_id")
	var g models.Hearts
	err := models.GetHearts(gr.Db, &g, gag_id)
	if err != nil {
		gr.checkErr(err, w)
		return
	}
	res := len(g)

	render.JSON(w, r, res)
}

func (gr *HeartRepo) Create(w http.ResponseWriter, r *http.Request) {
	gag_id, err := strconv.Atoi(chi.URLParam(r, "gag_id"))
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	user_id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	h := models.Heart{Gag_Id: gag_id, User_Id: user_id}
	err = models.CreateHeart(gr.Db, &h)
	if err != nil {
		gr.checkErr(err, w)
		return
	}

	render.JSON(w, r, "Successfully liked the gag")
}

func (gr *HeartRepo) Delete(w http.ResponseWriter, r *http.Request) {
	gag_id := chi.URLParam(r, "gag_id")
	user_id := chi.URLParam(r, "user_id")
	var h models.Heart
	err := models.DeleteHeart(gr.Db, &h, gag_id, user_id)
	if err != nil {
		gr.checkErr(err, w)
		return
	}

	http.Error(w, "Successfully dis-liked gag", http.StatusNotFound)
}
