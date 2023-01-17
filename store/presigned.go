package store

import (
	"context"

	"github.com/teacat/rushia/v3"
)

// Presigned
type Presigned struct {
	Token      string `rushia:"token" gorm:"column:token;primaryKey"`
	Type       string `rushia:"type" gorm:"column:type"`
	CreatedAt  int64  `rushia:"created_at" gorm:"column:created_at"`
	ExpiredAt  int64  `rushia:"expired_at" gorm:"column:expired_at"`
	FinishedAt int64  `rushia:"finished_at" gorm:"column:finished_at"`
}

// CreatePresigned
func (s *Store) CreatePresigned(c context.Context, data *Presigned) (err error) {
	err = s.Goshia.Exec(rushia.NewQuery("presigneds").Insert(data))
	return
}
