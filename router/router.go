package router

import (
	"context"
	"encoding/json"
	"errors"

	//	"log"
	"html/template"
	"net/http"
	"savemanager/database"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var db database.DataBase

func GeneralCtx(name string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ID := chi.URLParam(r, "ID")
			article, err := dbGet(name, ID)
			if err != nil {
				http.Error(w, http.StatusText(404), 404)
				return
			}
			ctx := context.WithValue(r.Context(), "article", article)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

var ProfileCtx = GeneralCtx("Profile")
var SaveDirectoryCtx = GeneralCtx("SaveDirectory")

// ///////////////////////////////////////////////////////////////////////////////////////////////////
func InitRouter() (*chi.Mux, error) {

	db := database.DataBase{Name: "main"}
	err := db.Connect()
	if err != nil {
		return nil, err
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.Get("/i", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("templates/profiles.html")
		if err != nil {
			println(err)
			println("Error with template !")
			return
		}
		var data ProfilesTpl
		profiles, err := db.GetAllProfiles()
		if err != nil {
			println(err.Error())
			println(profiles)
			println("error getting profiles !")
			return
		}
		length := len(*profiles)
		if length > 0 {
			data = ProfilesTpl{"Profile List", profiles}
		} else {

			data = ProfilesTpl{"No profiles found !", profiles}
		}

		t.Execute(w, data)
		//	w.Write([]byte(t))
	})
	r.Route("/api/v1", func(r chi.Router) {

		r.Route("/profile", func(r chi.Router) {
			r.Route("/{ID}", ProfileRoute)
		})
		r.Route("/save-directory", func(r chi.Router) {
			r.Route("/{ID}", SaveDirectoryRoute)
		})
	})
	return r, nil
}

type ProfilesTpl struct {
	Title       string
	ProfileList *[]*database.Profile
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////
func PluckAndSendFromContext(w http.ResponseWriter, r *http.Request) {
	jsonRes, err := json.Marshal(r.Context().Value("article"))
	if err != nil {
		println("Error : json could not be processed")
		return
	}
	w.Write(jsonRes)
}
func ProfileRoute(r chi.Router) {
	r.Use(ProfileCtx)
	r.Get("/", PluckAndSendFromContext)
}
func SaveDirectoryRoute(r chi.Router) {
	r.Use(SaveDirectoryCtx)
	r.Get("/", PluckAndSendFromContext)
}

func dbGet(name string, id string) (article interface{}, err error) {
	uint64Id, err := strconv.ParseUint(id, 10, 0)
	if err == nil {
		return nil, errors.New("Invalid ID")
	}
	uintId := uint(uint64Id)
	switch name {
	case "Profile":
		return db.GetProfileByID(uintId)
	case "SaveDirectory":
		return db.GetSaveDirectoryByID(uintId)
	default:
		return nil, errors.New("This db case does not exist")
	}
}

// TODO :
func dbPost(name string, profileOrDirectoryData interface{}) (article interface{}, err error) {
	//uint64Id, err := strconv.ParseUint(id, 10, 0)
	// if err == nil {
	// 	return nil, errors.New("Invalid ID")
	// }
	//uintId := uint(uint64Id)
	switch name {
	case "Profile":
		return db.CreateProfile(&database.Profile{})
	case "SaveDirectory":
		return db.CreateSaveDirectory(&database.SaveDirectory{})
	default:
		return nil, errors.New("This db case does not exist")
	}
}
