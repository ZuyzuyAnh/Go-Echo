package handler

import (
	"echo-demo/config"

	"gorm.io/gorm"
)

type (
	Handler struct {
		DB  *gorm.DB
		Cfg config.Config
	}
)

func NewHandler(DB *gorm.DB, cfg config.Config) *Handler {
	return &Handler{
		DB:  DB,
		Cfg: cfg,
	}
}
