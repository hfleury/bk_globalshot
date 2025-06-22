package service

import (
	"context"
)

type HealthService interface {
	Check(ctx context.Context) error
}

type dbHealthService struct {
	dbCheck func(context.Context) error
}

func NewDBHealthService(dbCheck func(context.Context) error) HealthService {
	return &dbHealthService{
		dbCheck: dbCheck,
	}
}

func (s *dbHealthService) Check(ctx context.Context) error {
	return s.dbCheck(ctx)
}
