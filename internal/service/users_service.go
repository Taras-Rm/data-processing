package service

import (
	"context"
	"fmt"
	"processor/internal/es"
	"processor/internal/infra/csv"
)

type UsersServiceI interface {
	ProcessAll(ctx context.Context) error
}

type usersService struct {
	usersStore es.UsersStore
	usersCsv   csv.UsersCsvI
}

func NewUsersService(usersStore es.UsersStore, usersCsv csv.UsersCsvI) UsersServiceI {
	return &usersService{
		usersStore,
		usersCsv,
	}
}

func (s *usersService) ProcessAll(ctx context.Context) error {
	users, err := s.usersStore.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("get all users: %w", err)
	}

	err = s.usersCsv.WriteAll(users)
	if err != nil {
		return fmt.Errorf("write all users: %w", err)
	}

	fmt.Println(len(users))

	return nil
}
