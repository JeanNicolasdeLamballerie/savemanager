package router

import (
	"context"
	"encoding/json"
	"errors"
	// "io/fs"
	"os"
	// "slices"

	// "fmt"
	"image/color"
	"path/filepath"
	"strings"
	//	"log"
	"html/template"
	"net/http"
	"savemanager/database"
	"savemanager/filemanager"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ncruces/zenity"
	"github.com/olahol/melody"
)

var db database.DataBase

type ServerSuccess struct {
	Resource string `json:"resource"`
	Status   string `json:"status"`

	Data UnknownData `json:"data"`
}

type UnknownData struct {
	Value any `json:"value"`
}

type SaveDirectoryInformations struct {
	Information map[string]string
	Tags        []string
}

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

	workDir, _ := os.Getwd()
	// TODO add backup db ?
	db := database.DataBase{Name: "main"}
	err := db.Connect()
	if err != nil {
		return nil, err
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(filepath.Join(workDir, "templates/profiles.html"))
		if err != nil {
			println(err.Error())
			println("Error with template !")
			return
		}
		var data ProfilesTpl
		profiles, err := db.GetAllProfiles()
		if err != nil {
			println(err.Error())
			println("error getting profiles !")
			return
		}
		length := len(*profiles)
		var title string
		if length > 0 {
			title = "Profile List"
		} else {
			title = "Your first setup..."
		}
		data = ProfilesTpl{title, profiles}
		t.Execute(w, data)
		//	w.Write([]byte(t))
	})
	r.Route("/api/v1", func(r chi.Router) {

		m := melody.New()
		/////////////////////////////////////////////////////////////////////////////////////////
		/////////////////////////////////////////////////////////////////////////////////////////
		//websocket API
		r.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
			m.HandleRequest(w, r)
		})
		m.HandleMessage(func(s *melody.Session, b []byte) {
			var msg WebsocketMessage
			err := json.Unmarshal(b, &msg)
			if err != nil {
				println("error :")
				println(err.Error())
			}

			switch msg.Resource {
			case "profile":
				shards := strings.Split(msg.Status, ":")
				switch shards[0] {
				case "init":
					if len(shards) > 1 {
						switch shards[1] {
						case "request-file-path":
							// todo put in savefile ?
							path, err := zenity.SelectFile(zenity.Directory(), zenity.Color(color.Black))
							if err != nil {
								println(err.Error())
								return
							}
							filemanager.CheckGameType(path)
							response := ServerSuccess{
								Resource: "profile",
								Status:   "success:file-path",
								Data: UnknownData{
									Value: path,
								},
							}
							broadcasted, err := json.Marshal(response)
							if err != nil {
								println(err.Error())
								return
							}
							m.Broadcast(broadcasted)
						case "register":
							m.Broadcast(b)
							if msg.Data.GamePath == "" || msg.Data.ProfileName == "" {
								println("ERROR : missing Path or Profile")
								//m.Broadcast("some error")
								return
							}
							newProfile := database.Profile{ProfileName: msg.Data.ProfileName, GamePath: msg.Data.GamePath}
							// TODO include sanity check ? e.g newProfile.verify(), maybe with a regex to check the path ?
							db.CreateProfile(&newProfile)

							stringifiedProfile, err := json.Marshal(newProfile)
							println(stringifiedProfile)
							if err != nil {
								println(err.Error())
								return
							}
							m.Broadcast(stringifiedProfile)
						}
					}

				}

			}
		})
		/////////////////////////////////////////////////////////////////////////////////////////
		/////////////////////////////////////////////////////////////////////////////////////////
		//rest API
		r.Route("/profile", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {

			})
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

type WebsocketMessage struct {
	Resource string // profile/else
	Status   string
	Data     struct {
		ProfileName string
		GamePath    string
	}
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

/// SOCKET CODE

// func handleSocket(w http.ResponseWriter, r *http.Request) {

// }
