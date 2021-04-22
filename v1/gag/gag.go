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
	"github.com/pamrulla/gagster-feed/helpers"
	"github.com/pamrulla/gagster-feed/models"
	"gorm.io/gorm"
)

type GetGagsResult struct {
	Path       string
	AthorId    string
	AuthorName string
	Width      int
	Height     int
	Price      float32
	Title      string
	Descrption string
	Tags       []string
	Hearts     int
}

type GagRepo struct {
	Db *gorm.DB
}

func New() *GagRepo {
	db := database.InitDb()
	db.AutoMigrate(models.Gag{})
	return &GagRepo{Db: db}
}

func (gr *GagRepo) checkErr(err error, w http.ResponseWriter, r *http.Request) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		helpers.NewError(w, r, "Gag not found", http.StatusNotFound)
	} else {
		helpers.NewError(w, r, "Something went wrong, please try again", http.StatusInternalServerError)
	}
}

func (gr *GagRepo) GetGags(w http.ResponseWriter, r *http.Request) {
	user_id := chi.URLParam(r, "user_id")
	var res models.Gags
	err := models.GetGags(gr.Db, &res, user_id)
	if err != nil {
		gr.checkErr(err, w, r)
		return
	}
	render.JSON(w, r, res)
}

func (gr *GagRepo) Create(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.NewError(w, r, "Invalid data sent", http.StatusBadRequest)
		return
	}
	var someData map[string]interface{}
	err = json.Unmarshal(req, &someData)
	if err != nil {
		helpers.NewError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	_, ok := someData["file"]
	if !ok {
		helpers.NewError(w, r, "File to upload is not present", http.StatusBadRequest)
		return
	}
	file, _ := someData["file"].(string)

	_, ok = someData["title"]
	if !ok {
		helpers.NewError(w, r, "Title is not present", http.StatusBadRequest)
		return
	}
	title, _ := someData["title"].(string)

	_, ok = someData["description"]
	if !ok {
		helpers.NewError(w, r, "Description is not present", http.StatusBadRequest)
		return
	}
	description, _ := someData["description"].(string)

	_, ok = someData["tags"]
	if !ok {
		helpers.NewError(w, r, "Tags are not present", http.StatusBadRequest)
		return
	}
	var tags []string
	for _, t := range someData["tags"].([]interface{}) {
		tags = append(tags, t.(string))
	}
	if err != nil {
		helpers.NewError(w, r, "failed to parse tags", http.StatusBadRequest)
		return
	}

	if _, ok := someData["price"]; !ok {
		helpers.NewError(w, r, "Price is not present", http.StatusBadRequest)
		return
	}

	u := models.Gag{}
	p, ok := someData["price"].(float64)
	if !ok {
		helpers.NewError(w, r, "Invalid price sent", http.StatusBadRequest)
		return
	}
	u.Price = float32(p)

	// err = json.Unmarshal(req, &u)
	// if err != nil {
	// 	helpers.NewError(w, r, "Invalid data sent", http.StatusBadRequest)
	// 	return
	// }

	user_id := chi.URLParam(r, "user_id")
	u.User_Id, err = strconv.Atoi(user_id)
	if err != nil {
		helpers.NewError(w, r, "Invalid User Id sent", http.StatusBadRequest)
		return
	}

	path, err := helpers.UploadCloudinary(file, user_id, tags)
	if err != nil {
		helpers.NewError(w, r, "Failed to upload image, please try again...", http.StatusInternalServerError)
		return
	}

	u.Path = path
	u.Title = title
	u.Description = description

	err = models.CreateGag(gr.Db, &u)
	if err != nil {
		gr.checkErr(err, w, r)
		return
	}
	render.JSON(w, r, "Successfully added new gag")
}

func (gr *GagRepo) Update(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.NewError(w, r, "Invalid data sent", http.StatusBadRequest)
		return
	}
	var u models.Gag
	err = json.Unmarshal(req, &u)
	if err != nil {
		helpers.NewError(w, r, "Invalid data sent", http.StatusBadRequest)
		return
	}

	err = models.UpdateGag(gr.Db, &u)
	if err != nil {
		gr.checkErr(err, w, r)
		return
	}
	render.JSON(w, r, u)
}

func (gr *GagRepo) Delete(w http.ResponseWriter, r *http.Request) {
	gag_id := chi.URLParam(r, "gag_id")

	var gag models.Gag
	err := models.DeleteGag(gr.Db, &gag, gag_id)
	if err != nil {
		gr.checkErr(err, w, r)
		return
	}

	render.JSON(w, r, "Successfully deleted gag")
}

func (gr *GagRepo) Get(w http.ResponseWriter, r *http.Request) {
	gag_id := chi.URLParam(r, "gag_id")
	var g models.Gag
	err := models.GetGag(gr.Db, &g, gag_id)
	if err != nil {
		gr.checkErr(err, w, r)
		return
	}
	render.JSON(w, r, g)
}

func (gr *GagRepo) Enable(w http.ResponseWriter, r *http.Request) {
	gag_id := chi.URLParam(r, "gag_id")

	var g models.Gag
	err := models.GetGag(gr.Db, &g, gag_id)
	if err != nil {
		gr.checkErr(err, w, r)
		return
	}
	g.IsEnabled = true
	err = models.UpdateGag(gr.Db, &g)
	if err != nil {
		gr.checkErr(err, w, r)
		return
	}
	render.JSON(w, r, "Successfully enabled gag")
}

func (gr *GagRepo) Disable(w http.ResponseWriter, r *http.Request) {
	gag_id := chi.URLParam(r, "gag_id")

	var g models.Gag
	err := models.GetGag(gr.Db, &g, gag_id)
	if err != nil {
		gr.checkErr(err, w, r)
		return
	}
	g.IsEnabled = false
	err = models.UpdateGag(gr.Db, &g)
	if err != nil {
		gr.checkErr(err, w, r)
		return
	}
	render.JSON(w, r, "Successfully disabled gag")
}

type GagsAsset struct {
	Id          int
	First_name  string
	Last_name   string
	User_id     int
	Path        string
	Price       float32
	Title       string
	Description string
	Hearts      int
	Width       int
	Height      int
	Tags        []string `gorm:"type:text"`
}
type GagsResult struct {
	Gags       []GagsAsset
	NextCursor string `json:"next_cursor"`
}

func (gr *GagRepo) Feed(w http.ResponseWriter, r *http.Request) {
	cld_gags, err := helpers.FeedCloudinary(r.URL.Query().Get("next_cursor"))
	if err != nil {
		helpers.NewError(w, r, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}

	if cld_gags.Assets == nil || len(cld_gags.Assets) == 0 {
		render.JSON(w, r, GagsResult{})
		return
	}

	pathSearch := ""
	arrLen := len(cld_gags.Assets)
	for i, g := range cld_gags.Assets {
		pathSearch += "path='" + g.PublicID + "'"
		if i != (arrLen - 1) {
			pathSearch += " OR "
		}
	}

	var result GagsResult
	gr.Db.Raw(`SELECT u.first_name, u.last_name, g.id, g.user_id, g.path, g.price, g.title, g.description, g.hearts  
		FROM gags g, users u
		WHERE (` + pathSearch + `) AND g.user_id = u.id ORDER BY g.created_at DESC`).Scan(&result.Gags)

	for i, _ := range result.Gags {
		for _, g := range cld_gags.Assets {
			if result.Gags[i].Path == g.PublicID {
				result.Gags[i].Width = g.Width
				result.Gags[i].Height = g.Height
				result.Gags[i].Tags = append(result.Gags[i].Tags, g.Tags...)
			}
		}
	}
	if len(result.Gags) < helpers.MAX_GAGS {
		result.NextCursor = ""
	} else {
		result.NextCursor = cld_gags.NextCursor
	}

	render.JSON(w, r, result)
}

func (gr *GagRepo) GetAuthorGags(w http.ResponseWriter, r *http.Request) {
	user_id := chi.URLParam(r, "user_id")
	cld_gags, err := helpers.GagsOfAuthor(user_id, r.URL.Query().Get("next_cursor"))

	if err != nil {
		helpers.NewError(w, r, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}

	if cld_gags.Assets == nil || len(cld_gags.Assets) == 0 {
		render.JSON(w, r, GagsResult{})
		return
	}

	pathSearch := ""
	arrLen := len(cld_gags.Assets)
	for i, g := range cld_gags.Assets {
		pathSearch += "path='" + g.PublicID + "'"
		if i != (arrLen - 1) {
			pathSearch += " OR "
		}
	}

	var result GagsResult
	gr.Db.Raw(`SELECT u.first_name, u.last_name, g.id, g.user_id, g.path, g.price, g.title, g.description, g.hearts  
		FROM gags g, users u
		WHERE (` + pathSearch + `) AND g.user_id = u.id ORDER BY g.created_at DESC`).Scan(&result.Gags)

	for i, _ := range result.Gags {
		for _, g := range cld_gags.Assets {
			if result.Gags[i].Path == g.PublicID {
				result.Gags[i].Width = g.Width
				result.Gags[i].Height = g.Height
				result.Gags[i].Tags = append(result.Gags[i].Tags, g.Tags...)
			}
		}
	}
	if len(result.Gags) < helpers.MAX_GAGS {
		result.NextCursor = ""
	} else {
		result.NextCursor = cld_gags.NextCursor
	}

	render.JSON(w, r, result)
}

func (gr *GagRepo) GagsWithTags(w http.ResponseWriter, r *http.Request) {
	cld_gags, err := helpers.GagsWithTags(r.URL.Query().Get("tags"), r.URL.Query().Get("next_cursor"))

	if err != nil {
		helpers.NewError(w, r, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}

	if cld_gags.Assets == nil || len(cld_gags.Assets) == 0 {
		render.JSON(w, r, GagsResult{})
		return
	}

	pathSearch := ""
	arrlen := len(cld_gags.Assets)
	for i, g := range cld_gags.Assets {
		pathSearch += "path='" + g.PublicID + "'"
		if i != (arrlen - 1) {
			pathSearch += " OR "
		}
	}

	var result GagsResult
	gr.Db.Raw(`SELECT u.first_name, u.last_name, g.id, g.user_id, g.path, g.price, g.title, g.description, g.hearts  
		FROM gags g, users u
		WHERE (` + pathSearch + `) AND g.user_id = u.id ORDER BY g.created_at DESC`).Scan(&result.Gags)

	for i, _ := range result.Gags {
		for _, g := range cld_gags.Assets {
			if result.Gags[i].Path == g.PublicID {
				result.Gags[i].Width = g.Width
				result.Gags[i].Height = g.Height
				result.Gags[i].Tags = append(result.Gags[i].Tags, g.Tags...)
			}
		}
	}
	if len(result.Gags) < helpers.MAX_GAGS {
		result.NextCursor = ""
	} else {
		result.NextCursor = cld_gags.NextCursor
	}

	render.JSON(w, r, result)
}
