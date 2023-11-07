package usecase

import "github.com/FackOff25/GoToTeamGradSuggester/internal/domain"

func (uc *UseCase) GetUser(uuid string) (*domain.User, error) {
	return uc.repo.GetUser()
}

func (uc *UseCase) ApplyUserReactionToPlace(uuid string, placeId string, reaction string) error {
	return nil
}
