package usecase

import "github.com/FackOff25/GoToTeamGradSuggester/internal/domain"

func (uc *UseCase) GetUser() (*domain.User, error) {
	return uc.repo.GetUser()
}
