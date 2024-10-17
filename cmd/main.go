package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"song-lib/internal/config"
	"song-lib/internal/user"
	"song-lib/internal/user/db"
	"song-lib/pkg/client"
	"song-lib/pkg/logging"
	"time"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/julienschmidt/httprouter"
	"github.com/pressly/goose/v3"
	"github.com/joho/godotenv"
	"os"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/gin-gonic/gin"
)

// @title Song-libary Swagger API 
// @version 1.0
// @description Swagger API for Golang Project Song-library
// @host localhost:10000
// @BasePath /
func main(){
	router := httprouter.New()

	router.GET("/swagger/*any", ginToHttprouter(ginSwagger.WrapHandler(swaggerFiles.Handler)))

	cfg := config.GetConfig()
	logger := logging.GetLogger()

	rep := databaseConnection(logger)
	
	handler := user.NewHandler(rep, logger)
	handler.Register(router)

	start(router, cfg, logger)
}

func databaseConnection(logger *logging.Logger) user.Repository{
	err := godotenv.Load("C:/Users/vadim/song-lib/.env")
	if err != nil {
        logger.Fatalf("Error loading .env file: %v", err)
    }

	client, err := postgresql.NewClient(context.Background(), 5)
	if err != nil {
		logger.Fatalf("Error: %v", err)
	}

	rep := db.NewRepository(client, logger)

	db, err := sql.Open("pgx", fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")))
	if err != nil {
		logger.Fatalf("Error of migration: %v", err)
	}
	defer db.Close()

	if err := goose.Up(db, "C:/Users/vadim/song-lib/migrations"); err != nil {
        logger.Fatal(err)
    }

	return rep
}

func start(router *httprouter.Router, cfg *config.Config, logger *logging.Logger) {
	logger.Info("start application")

	logger.Info("listen tcp")
	listener, listenErr := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler: router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}

func ginToHttprouter(h gin.HandlerFunc) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        c, _ := gin.CreateTestContext(w)
        c.Request = r
        h(c)
    }
}