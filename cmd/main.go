package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/db/psql"
	"github.com/hfleury/bk_globalshot/internal/handler"
	"github.com/hfleury/bk_globalshot/internal/router"
	"github.com/hfleury/bk_globalshot/internal/service"
	"github.com/hfleury/bk_globalshot/pkg/config"
	"github.com/hfleury/bk_globalshot/pkg/db"
	"github.com/hfleury/bk_globalshot/pkg/token"
)

func main() {
	cfg := config.LoadConfig()

	dbPsql, err := db.NewPsqlDb(cfg.DbDsn)
	if err != nil {
		panic(err)
	}

	maker, err := token.NewPasetoMaker(cfg.CfgToken.TokenKey)
	if err != nil {
		panic(err)
	}

	// Init repositories
	userRepo := psql.NewPostgresUserRepository(dbPsql)

	// Initi servies
	authService := service.NewAuthService(userRepo, maker)

	// Init handlers
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()
	// Set up CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router := router.NewRouter(r)
	router.SetupRouter(authHandler)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
