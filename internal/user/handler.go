package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"song-lib/internal/apperror"
	"song-lib/pkg/logging"
	"song-lib/pkg/utils"
	"github.com/gin-gonic/gin"
	_ "song-lib/docs"
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

func (h *handler) Register(router *gin.Engine) {
	router.GET("/search", apperror.Middleware(h.GetList))
	router.GET("/info", apperror.Middleware(h.GetInfo))
	router.DELETE("/remove", apperror.Middleware(h.Delete))
	router.PATCH("/update", apperror.Middleware(h.Update))
	router.POST("/new", apperror.Middleware(h.Create))
}

type SongShort struct {
    Group string `json:"group"`
    Song  string `json:"song"`
}


// GetListTags		godoc
// @Summary			Retrieves songs based on given parametres
// @Produce 		json
// @Param 			group query string false "group"
// @Param 			song query string false "song"
// @Success 		200 {object} user.Song
// @Failure			400 "Invalid input"
// @Failure			500 "Internal server error"
// @Router 			/search [get]
func (h *handler) GetList(c *gin.Context) error {
	var songs []Song
	var err error

	if len(c.Request.URL.Query()) == 0 {
		songs, err = h.repository.FindAll(context.Background())
		if err != nil {
			h.logger.Error(err)
			return apperror.DatabaseErr
		}
	} else {
		var s Song
		if c.Request.URL.Query().Get("group") != "" {
			s.Group = c.Request.URL.Query().Get("group")
		} else if c.Request.URL.Query().Get("song") != "" {
			s.Song = c.Request.URL.Query().Get("song")
		} else {
			return apperror.BadRequestErr
		}

		songs, err = h.repository.FindWithFilter(context.Background(), &s)
		if err != nil {
			h.logger.Error(err)
			return err
		}
	}

	var result []SongShort
    for _, s := range songs {
        result = append(result, SongShort{
            Group: s.Group,
            Song:  s.Song,
        })
    }

    c.JSON(http.StatusOK, result)

	return nil
}

// GetInfoTgs		godoc
// @Summary			Retrieves info about song based on given parametres
// @Produce			json
// @Param			group query string false "group"
// @Param			song query string false "song"
// @Success			200 {object} user.Song
// @Failure			400 "Invalid input"
// @Failure			500 "Internal server error"
// @Router			/info [get]
func (h *handler) GetInfo(ctx *gin.Context) error { 
	s := Song{Group: ctx.Request.URL.Query().Get("group"), Song: ctx.Request.URL.Query().Get("song")}
	if s.Group == "" || s.Song == "" {
		h.logger.Error(fmt.Errorf("bad request"))
		return apperror.BadRequestErr
	}

	song, err := h.repository.FindOne(context.TODO(), &s)
	if err != nil {
		h.logger.Error(err)
		return apperror.DatabaseErr
	}

	ctx.JSON(http.StatusOK, song)

	return nil
}

// DeleteTags		godoc
// @Summary			Deleting a song by its group and name
// @Param			group query string false "group"
// @Param			song query string false "song"
// @Success			204 "No content"
// @Failure			404 "User not found"
// @Failure			500 "Internal server error"
// @Router			/remove [delete]
func (h *handler) Delete(ctx *gin.Context) error {
	s := Song{Group: ctx.Request.URL.Query().Get("group"), Song: ctx.Request.URL.Query().Get("song")}

	err := h.repository.Delete(context.TODO(), &s)
	if err != nil {
		h.logger.Error(err)
		return apperror.DatabaseErr
	}

	ctx.Status(http.StatusOK)

	return nil
}

// UpdateTags		godoc
// @Summary			Update song information 
// @Param			group query string false "group"
// @Param			song query string false "song"
// @Success			200 "Data update"
// @Failure			404 "User not found"
// @Failure			500 "Internal server error"
// @Router			/update [patch]
func (h *handler) Update(ctx *gin.Context) error {

	if !repeatable.CheckContentType(ctx) {
		h.logger.Error(fmt.Errorf("content type is not json"))
		return apperror.BadRequestErr
	}

	fmt.Println(len(ctx.Request.URL.Query()))
	if len(ctx.Request.URL.Query()) == 2 {
		if ctx.Request.URL.Query().Get("group") != "" && ctx.Request.URL.Query().Get("song") != "" {
			s := Song{Group: ctx.Request.URL.Query().Get("group"), Song: ctx.Request.URL.Query().Get("song")}
			decoder := json.NewDecoder(ctx.Request.Body)
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

// CreateTags		godoc
// @Summary			Adding a new song to the library
// @Accept			json
// @Success			201 "Created"
// @Failure			400 "Invalid input"
// @Failure			500 "Internal server error"
// @Router			/new [post]
func (h *handler) Create(ctx *gin.Context) error {

	if !repeatable.CheckContentType(ctx) {
		h.logger.Error(fmt.Errorf("content type is not json"))
		return apperror.BadRequestErr
	}

	var s Song
	
	decoder := json.NewDecoder(ctx.Request.Body)
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
