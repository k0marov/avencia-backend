package service_helpers

import "github.com/AvenciaLab/avencia-backend/lib/core/db"

type DBService[INP, OUT any] func(db.DB, INP) (OUT, error)
type Service[INP, OUT any] func(INP) (OUT, error)

func NewDBTransService[INP any, OUT any](runT db.TransRunner, service DBService[INP, OUT]) Service[INP, OUT] {
	return func(i INP) (OUT, error) {
		var out OUT
		err := runT(func(d db.DB) error {
			res, err := service(d, i)
			out = res
			return err
		})
		return out, err
	}
}

type Nothing struct{}

func NewNoResultService[INP any](service func(INP) error) Service[INP, Nothing] {
	return func(i INP) (Nothing, error) {
		return Nothing{}, service(i)
	}
}

func NewDBNoResultService[INP any](runT db.TransRunner, service func(db.DB, INP) error) Service[INP, Nothing] {
  return func(i INP) (Nothing, error) {
    err := runT(func(db db.DB) error {
       return service(db, i) 
    })
    return Nothing{}, err
  }
}

// TODO: maybe rename all package names to camel case since it is idiomatic
