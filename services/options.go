package services

import "samwang0723/jarvis/db/dal/idal"

type Option func(o *serviceImpl)

func WithDAL(dal idal.IDAL) Option {
	return func(i *serviceImpl) {
		i.dal = dal
	}
}
