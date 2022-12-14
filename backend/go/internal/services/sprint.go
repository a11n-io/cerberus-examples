package services

import (
	"cerberus-examples/internal/database"
	"cerberus-examples/internal/repositories"
	"context"
	"fmt"
)

type SprintService interface {
	Create(ctx context.Context, projectId, goal string) (repositories.Sprint, error)
	FindByProject(ctx context.Context, projectId string) ([]repositories.Sprint, error)
	Get(ctx context.Context, sprintId string) (repositories.Sprint, error)
	Start(ctx context.Context, sprintId string) (repositories.Sprint, error)
	End(ctx context.Context, sprintId string) (repositories.Sprint, error)
}

type sprintService struct {
	txProvider database.TxProvider
	repo       repositories.SprintRepo
}

func NewSprintService(
	txProvider database.TxProvider,
	repo repositories.SprintRepo) SprintService {
	return &sprintService{
		txProvider: txProvider,
		repo:       repo,
	}
}

func (s *sprintService) Create(ctx context.Context, projectId, goal string) (repositories.Sprint, error) {

	userId := ctx.Value("userId")
	if userId == nil {
		return repositories.Sprint{}, fmt.Errorf("no userId")
	}

	tx, err := s.txProvider.GetTransaction()
	if err != nil {
		return repositories.Sprint{}, err
	}

	sprint, err := s.repo.Create(projectId, goal, tx)
	if err != nil {
		if rbe := tx.Rollback(); rbe != nil {
			err = fmt.Errorf("rollback error (%v) after %w", rbe, err)
		}
		return repositories.Sprint{}, err
	}

	return sprint, tx.Commit()
}

func (s *sprintService) FindByProject(ctx context.Context, projectId string) ([]repositories.Sprint, error) {
	return s.repo.FindByProject(projectId)
}

func (s *sprintService) Get(ctx context.Context, sprintId string) (repositories.Sprint, error) {
	return s.repo.Get(sprintId, nil)
}

func (s *sprintService) Start(ctx context.Context, sprintId string) (repositories.Sprint, error) {
	_, err := s.repo.Start(sprintId)
	if err != nil {
		return repositories.Sprint{}, err
	}
	return s.repo.Get(sprintId, nil)
}

func (s *sprintService) End(ctx context.Context, sprintId string) (repositories.Sprint, error) {
	_, err := s.repo.End(sprintId)
	if err != nil {
		return repositories.Sprint{}, err
	}
	return s.repo.Get(sprintId, nil)
}
