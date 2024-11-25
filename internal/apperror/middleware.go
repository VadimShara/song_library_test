package apperror

import (
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
)

type apphandler func(c *gin.Context) error

func Middleware(h apphandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var appErr *AppError
		err := h(ctx)
		if err != nil {
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrNotFound) {
					ctx.JSON(http.StatusNotFound, gin.H{"error": ErrNotFound.Marshal()})
					return
				} else if errors.Is(err, DatabaseErr) {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": DatabaseErr.Marshal()})
					return
				}

				err = err.(*AppError)
				ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrNotFound.Marshal()})
				return
			}

			ctx.JSON(http.StatusTeapot, gin.H{"error": systemError(err).Marshal()})
		}
	}
}