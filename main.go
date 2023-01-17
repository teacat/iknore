package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/teacat/ginrpc"
	"github.com/teacatx/iknore/config"
	"github.com/teacatx/iknore/handler"
	"github.com/teacatx/iknore/service/file"
	"github.com/teacatx/iknore/store"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "start",
				Usage: "start the iknore file server",
				Action: func(c *cli.Context) error {
					router(c)
					return nil
				},
			},
			{
				Name:  "key",
				Usage: "generate the auth key",
				Action: func(*cli.Context) error {
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "remote",
				Aliases: []string{"r"},
				Value:   "local",
				Usage:   "available remote: s3, local (by using s3 please set AWS_ envs)",
			},
			&cli.StringFlag{
				Name:    "directory",
				Aliases: []string{"d"},
				Value:   "./files",
				Usage:   "directory for local remote",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func router(c *cli.Context) {
	r := gin.Default()
	conf := config.LoadConfig()

	//
	var fileSvc file.Handler
	switch c.String("remote") {
	case "s3":
		fileSvc = file.NewS3()
	case "local":
		fileSvc = file.NewLocal(c.String("directory"))
	default:
		log.Fatalln("no available remote")
	}
	//
	f := NewFactory(store.NewConn(), fileSvc, conf)

	//
	r.GET("/images/:type/:name", f.New(func(h *handler.Handler) gin.HandlerFunc {
		return func(c *gin.Context) {
			handler.ToContext(c, h)
			handler.GetImage(c)
		}
	}))
	//
	r.POST("/upload_presigned_image", f.New(func(h *handler.Handler) gin.HandlerFunc {
		return ginrpc.NewForm(h.UploadImage)
	}))
	//
	r.POST("/api/upload_image", f.New(func(h *handler.Handler) gin.HandlerFunc {
		return ginrpc.NewForm(h.UploadImage)
	}))
	//
	r.POST("/api/delete_image", f.New(func(h *handler.Handler) gin.HandlerFunc {
		return ginrpc.New(h.DeleteImage)
	}))
	//
	r.POST("/api/validate_image", f.New(func(h *handler.Handler) gin.HandlerFunc {
		return ginrpc.New(h.ValidateImage)
	}))
	//
	r.POST("/api/create_image_pointer", f.New(func(h *handler.Handler) gin.HandlerFunc {
		return ginrpc.New(h.CreateImagePointer)
	}))
	//
	r.POST("/api/delete_image_pointer", f.New(func(h *handler.Handler) gin.HandlerFunc {
		return ginrpc.New(h.DeleteImagePointer)
	}))
	//
	r.POST("/api/create_presigned_url", f.New(func(h *handler.Handler) gin.HandlerFunc {
		return ginrpc.New(h.DeleteImagePointer)
	}))
	//
	r.Run()
}
