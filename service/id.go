package service

import uuid "github.com/satori/go.uuid"

type IDService struct {
}

func NewIDService() *IDService {
	return &IDService{}
}

// NewUserID
func (i *IDService) NewImageID() string {
	return uuid.NewV4().String()
}
