package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/teacatx/iknore/config"
	"github.com/teacatx/iknore/service"
	"github.com/teacatx/iknore/service/file"
	"github.com/teacatx/iknore/store"
)

type Handler struct {
	GlobalStore  *store.Store
	Store        *store.Store
	IDService    *service.IDService
	FileService  file.Handler
	ImageService *service.ImageService
	Config       *config.Config
}

const (
	KeyHandler = "handler"
)

func FromContext(c context.Context) *Handler {
	return c.Value(KeyHandler).(*Handler)
}

func ToContext(c *gin.Context, handler *Handler) {
	c.Set(KeyHandler, handler)
}
