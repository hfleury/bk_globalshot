package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/handler"
	"github.com/hfleury/bk_globalshot/internal/repository/psql"
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

	pasetoMaker, err := token.NewPasetoMaker(cfg.CfgToken.TokenKey)
	if err != nil {
		panic(err)
	}

	// Init repositories
	userRepo := psql.NewPostgresUserRepository(dbPsql)
	companyRepo := psql.NewPostgresCompanyRepository(dbPsql)
	roomRepo := psql.NewPostgresRoomRepository(dbPsql)
	siteRepo := psql.NewSiteRepository(dbPsql)
	unitRepo := psql.NewUnitRepository(dbPsql)

	// Initi servies
	authService := service.NewAuthService(userRepo, pasetoMaker, &cfg.CfgToken)
	companyService := service.NewCompanyService(dbPsql, companyRepo, userRepo)
	roomService := service.NewRoomService(roomRepo)
	siteService := service.NewSiteService(dbPsql, siteRepo)
	unitService := service.NewUnitService(dbPsql, unitRepo)
	userService := service.NewUserService(userRepo)
	dbHealthService := service.NewDBHealthService(func(ctx context.Context) error {
		return dbPsql.PingContext(ctx)
	})

	// Init handlers
	authHandler := handler.NewAuthHandler(authService)
	companyHandler := handler.NewCompanyHandler(companyService)
	roomHandler := handler.NewRoomHandler(roomService)
	siteHandler := handler.NewSiteHandler(siteService)
	unitHandler := handler.NewUnitHandler(unitService)
	userHandler := handler.NewUserHandler(userService)
	healthHandler := handler.NewHealthHandler(dbHealthService)

	r := gin.Default()

	router := router.NewRouter(r)
	router.SetupRouter(authHandler, healthHandler, companyHandler, roomHandler, siteHandler, unitHandler, userHandler, pasetoMaker)

	port := cfg.ServerPort
	if port == "" {
		port = "8080" // default port
	}

	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}
