package heart

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/pamrulla/gagster-feed/models"
)

var hearts models.Hearts

func Init() {}

func getParamFromRequest(r *http.Request, key string) int {
	vars := chi.URLParam(r, key)
	val, _ := strconv.Atoi(vars)
	return val
}

func GetHearts(w http.ResponseWriter, r *http.Request) {
	gag_id := getParamFromRequest(r, "gag_id")
	res := 0
	for _, g := range hearts {
		if g.Gag_Id == gag_id {
			fmt.Println(res)
			res++
		}
	}
	render.JSON(w, r, res)
}

func Create(w http.ResponseWriter, r *http.Request) {
	gag_id := getParamFromRequest(r, "gag_id")
	user_id := getParamFromRequest(r, "user_id")

	h := models.Heart{Gag_Id: gag_id, User_Id: user_id, Liked_On: time.Now().UTC()}
	hearts = append(hearts, h)
	render.JSON(w, r, "Successfully liked the gag")
}

func Delete(w http.ResponseWriter, r *http.Request) {
	gag_id := getParamFromRequest(r, "gag_id")
	user_id := getParamFromRequest(r, "user_id")
	isFound := false

	for i, a := range hearts {
		if a.Gag_Id == gag_id && user_id == a.User_Id {
			hearts = append(hearts[:i], hearts[i+1:]...)
			isFound = true
			break
		}
	}
	if isFound {
		render.JSON(w, r, "Successfully dis-liked the gag")
	} else {
		http.Error(w, "Gag Like not found", http.StatusNotFound)
	}
}
