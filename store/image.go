package store

import (
	"context"

	"github.com/teacat/rushia/v3"
)

type CoverMode string

const (
	CoverModeNone      = ""
	CoverModeSmart     = "smart"
	CoverModeNorthWest = "top_left"
	CoverModeNorth     = "top"
	CoverModeNorthEast = "top_right"
	CoverModeWest      = "left"
	CoverModeCenter    = "center"
	CoverModeEast      = "right"
	CoverModeSouthWest = "bottom_left"
	CoverModeSouth     = "bottom"
	CoverModeSouthEast = "bottom_right"
	CoverModeContain   = "contain"
)

type ImageArguments struct {
	Type              string
	Extension         string
	Format            string
	Name              string
	Size              string
	Width             int
	Height            int
	Filename          string
	ContentType       string
	CoverMode         CoverMode
	BackgroundColor   string
	IgnoreAspectRatio bool
}

// Image
type Image struct {
	ID        string `rushia:"id" gorm:"column:id;primaryKey"`
	Type      string `rushia:"type" gorm:"column:type"`
	OwnerName string `rushia:"owner_name" gorm:"column:owner_name"`
	Width     int    `rushia:"width" gorm:"column:width"`
	Height    int    `rushia:"height" gorm:"column:height"`
	CreatedAt int64  `rushia:"created_at" gorm:"column:created_at"`
}

// CreateImage
func (s *Store) CreateImage(c context.Context, data *Image) (err error) {
	err = s.Goshia.Exec(rushia.NewQuery("images").Insert(data))
	return
}

// ValidateImage
func (s *Store) ValidateImage(c context.Context, id, typ, ownerName string) (exists bool, err error) {
	err = s.Goshia.Query(rushia.NewQuery("images").
		Where("type = ?", typ).
		Where("owner_name = ?", ownerName).
		Where("id = ?", id).Exists(), &exists)
	return
}
