package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teacat/goshia/v3"
	"github.com/teacatx/iknore/config"
	"github.com/teacatx/iknore/handler"
	"github.com/teacatx/iknore/service"
	"github.com/teacatx/iknore/service/file"
	"github.com/teacatx/iknore/store"
)

type Factory struct {
	Conn         *store.Conn
	IDService    *service.IDService
	FileService  file.Handler
	ImageService *service.ImageService
	Config       *config.Config
}

// Factory
func NewFactory(conn *store.Conn, fileSvc file.Handler, conf *config.Config) *Factory {
	idsvc := service.NewIDService()
	imgsvc := service.NewImageService(conf)
	return &Factory{
		Conn:         conn,
		IDService:    idsvc,
		FileService:  fileSvc,
		ImageService: imgsvc,
		Config:       conf,
	}
}

func (f *Factory) New(h func(*handler.Handler) gin.HandlerFunc) gin.HandlerFunc {
	newGoshia := goshia.New(f.Conn.DB)
	newStore := store.New(newGoshia)

	return func(c *gin.Context) {
		//
		err := newGoshia.Transaction(func(tx *goshia.Goshia) error {
			sessionStore := store.New(tx)
			//
			h(&handler.Handler{
				GlobalStore:  newStore,
				Store:        sessionStore,
				IDService:    f.IDService,
				FileService:  f.FileService,
				ImageService: f.ImageService,
				Config:       f.Config,
			})(c)
			//
			c.Next()
			//
			for _, v := range c.Errors {
				return v
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
	}
}
