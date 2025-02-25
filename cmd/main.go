package main

import (
	"echo-demo/config"
	"echo-demo/internal/handler"
	"echo-demo/internal/model"

	"fmt"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	e := echo.New()

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Println(cfg.DB_DSN)

	e.Logger.SetLevel(log.ERROR)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.JWTSecret),
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/login" || c.Path() == "/register" {
				return true
			}

			return false
		},
	}))

	db, err := initializeDB(cfg.DB_DSN)
	if err != nil {
		e.Logger.Fatal(err)
	}

	db.AutoMigrate(&model.User{})

	h := handler.NewHandler(db, *cfg)

	e.POST("/register", h.Signup)
	e.POST("/login", h.Login)

	e.Logger.Fatal(e.Start(":1323"))
}

func initializeDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&model.User{}, &model.Role{}, &model.Permission{}, &model.RolePermission{}, &model.UserRole{},
		&model.Movie{}, &model.Theater{}, &model.Seat{}, &model.Payment{}, &model.PaymentDetail{},
		&model.MovieCategory{}, &model.MovieCategoryMapping{},
	)

	return db, nil
}
