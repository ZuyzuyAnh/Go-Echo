package main

import (
	"echo-demo/config"
	"echo-demo/internal/controller"
	"echo-demo/internal/repository"
	"echo-demo/internal/service"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {

	e := echo.New()

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Println(cfg.DB_DSN)

	e.Logger.SetLevel(log.ERROR)
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.JWTSecret),
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/users/login" || c.Path() == "/users/register" {
				return true
			}

			return false
		},
	}))

	db, err := initializeDB(cfg.DB_DSN)
	if err != nil {
		e.Logger.Fatal(err)
	}

	setupRoutes(e, db, initLogger(), cfg.JWTSecret)

	e.Logger.Fatal(e.Start(":1323"))
}

func initLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger
}

func initializeDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect db err: %v", err)
	}

	return db, nil
}

func setupRoutes(e *echo.Echo, db *sqlx.DB, logger *zap.Logger, secret string) {
	ur := repository.NewUserRepository(db)
	rr := repository.NewRoleRepository(db)
	mr := repository.NewMovieRepository(db)

	us := service.NewUserService(ur, rr)
	ms := service.NewMovieService(mr)

	mc := controller.NewMovieController(ms, logger)
	uc := controller.NewUserController(us, logger, secret)

	users := e.Group("/users")
	{
		users.POST("/login", uc.Login)
		users.POST("/register", uc.Register)
		users.GET("/profile", uc.GetProfile)
	}

	movies := e.Group("/movies")
	{
		movies.POST("", mc.CreateMovie)
	}
}
