package repeatable

import(
	"github.com/gin-gonic/gin"
	"time"
)

func DoWithTries(fn func() error, attemtps int, delay time.Duration) (err error) {
	for attemtps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemtps--

			continue
		}

		return nil
	}

	return
}

func CheckContentType(ctx *gin.Context) bool {
	contentType := ctx.Request.Header.Get("Content-Type")
	return contentType == "application/json"
}