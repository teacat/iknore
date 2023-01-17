package handler

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teacatx/iknore/pb"
	"github.com/teacatx/iknore/service"
	"github.com/teacatx/iknore/store"
)

// UploadImage
func (h *Handler) UploadImage(c *gin.Context, input *pb.UploadImageRequest) (*pb.UploadImageResponse, error) {
	//
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		return nil, err
	}
	//
	uploadedFile, err := input.File.Open()
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(tmpFile, uploadedFile); err != nil {
		return nil, err
	}
	if err := tmpFile.Close(); err != nil {
		return nil, err
	}
	//
	width, height, err := service.Compress(c, tmpFile.Name())
	if err != nil {
		return nil, err
	}
	//
	image := &store.Image{
		ID:        h.IDService.NewImageID(),
		Type:      input.Type,
		OwnerName: input.OwnerName,
		Width:     width,
		Height:    height,
		CreatedAt: time.Now().Unix(),
	}
	if err := h.Store.CreateImage(c, image); err != nil {
		return nil, err
	}
	//
	b, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}
	//
	if err := h.FileService.Upload(c, fmt.Sprintf("/%s/%s", input.Type, image.ID), b); err != nil {
		return nil, err
	}
	return &pb.UploadImageResponse{
		ID: image.ID,
	}, nil
}
