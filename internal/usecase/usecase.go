package usecase

import (
	"context"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/repository"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
)

type UseCase struct {
	repo *repository.Repo
	ctx  context.Context
}

func New(r repository.Repo, ctx context.Context) UsecaseInterface {
	return &UseCase{repo: &r, ctx: ctx}
}

func (uc *UseCase) GetUserPlaces() (*domain.User, error) {
	return uc.repo.GetUserPlaces()
}
