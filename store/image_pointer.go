package store

import (
	"context"

	"github.com/teacat/rushia/v3"
)

// ImagePointer
type Pointer struct {
	ImageID string `rushia:"image_id" gorm:"column:image_id;primaryKey"`
	Name    string `rushia:"name" gorm:"column:name;primaryKey"`
}

// CreateImagePointer
func (s *Store) CreateImagePointer(c context.Context, data *Image) (err error) {
	err = s.Goshia.Exec(rushia.NewQuery("image_pointers").Insert(data))
	return
}
