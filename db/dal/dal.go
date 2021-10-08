package dal

import (
	"samwang0723/jarvis/db/dal/idal"

	"gorm.io/gorm"
)

type dalImpl struct {
	db *gorm.DB
}

func New(opts ...Option) idal.IDAL {
	impl := &dalImpl{}
	for _, opt := range opts {
		opt(impl)
	}
	return impl
}
