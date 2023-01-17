package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/teacatx/iknore/pb"
)

// ValidateImage
func (h *Handler) ValidateImage(c *gin.Context, input *pb.ValidateImageRequest) (*pb.ValidateImageResponse, error) {
	//
	exists, err := h.Store.ValidateImage(c, input.ID, input.Type, input.OwnerName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &pb.ValidateImageResponse{
			IsValid: false,
		}, nil
	}
	return &pb.ValidateImageResponse{
		IsValid: true,
	}, nil
}
