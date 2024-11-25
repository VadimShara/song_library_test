package main

import (
	_ "song-lib/docs"
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
	router := gin.Default()

	cfg := config.GetConfig()
	logger := logging.GetLogger()
	
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

func start(router *gin.Engine, cfg *config.Config, logger *logging.Logger) {
	logger.Info("start application")

	rep := databaseConnection(logger)
	logger.Info("database successfully connected and migrated")

	handler := user.NewHandler(rep, logger)
	handler.Register(router)

	url := ginSwagger.URL("http://localhost:10000/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

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