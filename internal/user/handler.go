package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"song-lib/pkg/utils"
	"song-lib/internal/apperror"
	"song-lib/pkg/logging"
	"github.com/julienschmidt/httprouter"
	// swaggerFiles "github.com/swaggo/files"
	// "github.com/swaggo/gin-swagger"
	// "github.com/gin-gonic/gin"
)

type handler struct {
	logger *logging.Logger
	repository Repository
}

func NewHandler(repository Repository, logger *logging.Logger) *handler {
	return &handler{
		repository : repository,
		logger : logger,
	}
}

// func ginToHttprouter(h gin.HandlerFunc) httprouter.Handle {
//     return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//         c, _ := gin.CreateTestContext(w)
//         c.Request = r
//         h(c)
//     }
// }

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/search", apperror.Middleware(h.GetList))
	router.HandlerFunc(http.MethodGet, "/info", apperror.Middleware(h.GetInfo))
	router.HandlerFunc(http.MethodDelete, "/remove", apperror.Middleware(h.Delete))
	router.HandlerFunc(http.MethodPatch, "/update", apperror.Middleware(h.Update))
	router.HandlerFunc(http.MethodPost, "/new", apperror.Middleware(h.Create))
	//router.GET("/swagger/*any", ginToHttprouter(ginSwagger.WrapHandler(swaggerFiles.Handler)))
}

// GetList godoc
// @Summary Retrieves songs based on given parametres
// @Produce json
// @Param group query string false "group"
// @Param song query string false "song"
// @Success 200 {object} model.Song
// @Router /search [get]

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	var songs []Song
	var err error

	if len(r.URL.Query()) == 0 {
		songs, err = h.repository.FindAll(context.Background())
		if err != nil {
			h.logger.Error(err)
			return apperror.DatabaseErr
		}
	} else {
		var s Song
		if r.URL.Query().Get("group") != "" {
			s.Group = r.URL.Query().Get("group")
		} else if r.URL.Query().Get("song") != "" {
			s.Song = r.URL.Query().Get("song")
		} else {
			return fmt.Errorf("error request")
		}

		songs, err = h.repository.FindWithFilter(context.Background(), &s)
		if err != nil {
			h.logger.Error(err)
			return err
		}
	}

	for _, s := range songs {
			res, err := json.Marshal(s)
			if err != nil {
				h.logger.Error(err)
				return apperror.MarshalJSONErr
			}
			w.Header().Set("Content-Type", "application/json")
			//w.Write([]byte(fmt.Sprintf("Group: %s Song: %s\n", s.Group, s.Song)))
			w.Write([]byte(res))
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

// GetInfo godoc
// @Summary Retrieves info about song based on given parametres
// @Produce json
// @Param group query string false "group"
// @Param song query string false "song"
// @Success 200 {object} model.Song
// @Router /info [get]

func (h *handler) GetInfo(w http.ResponseWriter, r *http.Request) error { 
	s := Song{Group: r.URL.Query().Get("group"), Song: r.URL.Query().Get("song")}

	song, err := h.repository.FindOne(context.TODO(), &s)
	if err != nil {
		h.logger.Error(err)
		return apperror.DatabaseErr
	}

	res, err := json.Marshal(song)
	if err != nil {
		h.logger.Error(err)
		return apperror.MarshalJSONErr
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(res))
	//w.Write([]byte(fmt.Sprintf("Release date: %s\n, Text: %s\n, Link: %s\n", song.ReleaseDate, song.Text, song.Link)))
	w.WriteHeader(http.StatusOK)
	return nil
}

// Delete godoc
// @Summary Deleting a song by its group and name
// @Param group query string false "group"
// @Param song query string false "song"
// @Success 204 "No content"
// @Failure 404 "User not found"
// @Failure 500 "Internal server error"
// @Router /remove [delete]

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) error {
	s := Song{Group: r.URL.Query().Get("group"), Song: r.URL.Query().Get("song")}

	err := h.repository.Delete(context.TODO(), &s)
	if err != nil {
		h.logger.Error(err)
		return apperror.DatabaseErr
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

// Update godoc
// @Summary Update song information 
// @Param group query string false "group"
// @Param song query string false "song"
// @Success 200 "No content"
// @Failure 404 "User not found"
// @Failure 500 "Internal server error"
// @Router /udate [delete]

func (h *handler) Update(w http.ResponseWriter, r *http.Request) error {

	if !repeatable.CheckContentType(r) {
		h.logger.Error(fmt.Errorf("content type is not json"))
		return apperror.BadRequestErr
	}

	fmt.Println(len(r.URL.Query()))
	if len(r.URL.Query()) == 2 {
		if r.URL.Query().Get("group") != "" && r.URL.Query().Get("song") != "" {
			s := Song{Group: r.URL.Query().Get("group"), Song: r.URL.Query().Get("song")}
			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()
			err := decoder.Decode(&s)
			if err != nil {
				h.logger.Error(err)
				return apperror.DecodingErr
			}
			err = h.repository.Update(context.Background(), &s)
			if err != nil {
				h.logger.Error(err)
				return apperror.DatabaseErr
			}

		} else {
			return apperror.BadRequestErr
		}
	} else {
		return apperror.BadRequestErr
	}

	return nil
}

// Create godoc
// @Summary Adding a new song to the library
// @Accept json
// @Success 201 "No content"
// @Failure 400 "Invalid input"
// @Failure 500 "Internal server error"
// @Router /new [post]

func (h *handler) Create(w http.ResponseWriter, r *http.Request) error {

	if !repeatable.CheckContentType(r) {
		h.logger.Error(fmt.Errorf("content type is not json"))
		return apperror.BadRequestErr
	}

	var s Song
	
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&s)
	if err != nil {
		h.logger.Error(err)
		return apperror.DecodingErr
	}

	err = h.repository.Create(context.TODO(), &s)
	if err != nil {
		return apperror.DatabaseErr
	}

	return nil
}