package usecase

import (
	"github.com/jackc/pgx/v5"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
)

func (uc *UseCase) AddUser(id string) error {
	_, err := uc.repo.GetUser(id)

	if err != nil {
		if err == pgx.ErrNoRows {
			err := uc.repo.AddUser(id)
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	return nil
}

func (uc *UseCase) GetUser(uuid string) (*domain.User, error) {
	return uc.repo.GetUser(uuid)
}

func (uc *UseCase) ApplyUserReactionToPlace(uuid string, placeId string, reaction string) error {
	return nil
}