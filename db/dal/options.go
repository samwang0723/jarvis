package dal

import "gorm.io/gorm"

type Option func(o *dalImpl)

func WithDB(db *gorm.DB) Option {
	return func(s *dalImpl) {
		s.db = db
	}
}
