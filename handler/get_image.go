package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teacatx/iknore/store"
)

func VariantFilename(args *store.ImageArguments) string {
	path := fmt.Sprintf("/%s/%s-", args.Type, args.Name)
	if args.Width != 0 {
		path += fmt.Sprintf("w_%d-", args.Width)
	}
	if args.Height != 0 {
		path += fmt.Sprintf("h_%d-", args.Height)
	}
	if args.CoverMode != store.CoverModeNone {
		path += fmt.Sprintf("c_%s-", args.CoverMode)
	}
	if args.BackgroundColor != "" {
		path += fmt.Sprintf("bc_%s-", args.BackgroundColor)
	}

	path += fmt.Sprintf("t_%d-", time.Now().Unix())

	return strings.TrimSuffix(path, "-") + args.Extension
}

func SuffixToContentType(v string) string {
	switch v {
	case ".png":
		return "image/png"
	case ".bmp":
		return "image/bmp"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".avif":
		return "image/avif"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

func SuffixToFormat(v string) string {
	switch v {
	case ".png":
		return "png"
	case ".bmp":
		return "bmp"
	case ".jpg", ".jpeg":
		return "jpg"
	case ".avif":
		return "avif"
	case ".gif":
		return "gif"
	case ".webp":
		return "webp"
	default:
		return ""
	}
}

// GetImageArguments
func GetImageArguments(c *gin.Context) *store.ImageArguments {
	//
	typ := c.Params.ByName("type")
	ext := filepath.Ext(c.Params.ByName("name"))
	name := strings.TrimSuffix(filepath.Base(c.Params.ByName("name")), ext)
	format := SuffixToFormat(ext)
	//
	width, _ := strconv.Atoi(c.Query("w"))
	height, _ := strconv.Atoi(c.Query("h"))
	bgColor := c.Query("bc")
	size := c.Query("s")
	ignoreAspectRatio := strings.ToLower(c.Query("iar")) == "true"
	//
	return &store.ImageArguments{
		Type:              typ,
		Extension:         ext,
		Format:            format,
		Name:              name,
		Width:             width,
		Height:            height,
		Size:              size,
		ContentType:       SuffixToContentType(ext),
		CoverMode:         store.CoverMode(c.Query("c")),
		BackgroundColor:   bgColor,
		IgnoreAspectRatio: ignoreAspectRatio,
	}
}

func GetImage(c *gin.Context) {
	h := FromContext(c)
	//
	args := GetImageArguments(c)
	if args.Size != "" && (args.Width != 0 || args.Height != 0) {
		c.AbortWithError(http.StatusInternalServerError, errors.New("cannot size and w/h same time"))
		return
	}
	if args.Size != "" {
		args.Width, args.Height = h.ImageService.AliasToSize(args)

		if args.Width == 0 && args.Height == 0 {
			c.AbortWithError(http.StatusInternalServerError, errors.New("wrong size?"))
			return
		}
	}

	filename := VariantFilename(args)
	//
	c.Writer.Header().Set("Content-Type", args.ContentType)

	if !h.ImageService.CheckValidFormat(args) {
		c.AbortWithError(http.StatusInternalServerError, errors.New("not setted ext"))
		return
	}

	err := h.FileService.StreamDownload(c, filename, c.Writer) // delete all variations if link change
	if err == nil {

		return
	}

	//
	if !h.FileService.IsNotExist(err) {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !h.ImageService.CheckValidBackgroundColor(args) {
		c.AbortWithError(http.StatusInternalServerError, errors.New("not valid bg color"))
		return
	}
	if !h.ImageService.CheckValidCover(args) {
		c.AbortWithError(http.StatusInternalServerError, errors.New("not valid cover"))
		return
	}

	if !h.ImageService.CheckVaildSize(args) {
		c.AbortWithError(http.StatusInternalServerError, errors.New("not valid size preset"))
		return
	}

	// check available formats limit

	//
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// 下載原圖
	err = h.FileService.StreamDownload(c, fmt.Sprintf("/%s/%s", args.Type, args.Name), tmpFile)
	if err != nil {
		if !h.FileService.IsNotExist(err) {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// 如果找不到原圖就從 ImageService 拿這個 type 對應的 Placeholders。
		b, err := h.ImageService.GetPlaceholder(args)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, args.ContentType, b)
		return

	}
	//
	if err := tmpFile.Close(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	//
	if err := h.ImageService.MakeVariant(c, tmpFile.Name(), args); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	b, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	go func() {
		if err := h.FileService.Upload(c, filename, b); err != nil {
			//c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}()

	c.Data(http.StatusOK, args.ContentType, b)
	return
}
