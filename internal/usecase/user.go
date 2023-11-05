package usecase

import (
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (uc *UseCase) GetUser(id string) (*domain.User, error) {
	dbUser, err := uc.repo.GetUser(id)

	if err != nil {
		if err == pgx.ErrNoRows {
			err := uc.repo.AddUser(id)
			if err != nil {
				return nil, err
			}
			return &domain.User{Id: uuid.MustParse(id)}, nil

		} else {
			return nil, err
		}
	}

	return dbUser, nil
}
