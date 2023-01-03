package services

import (
	"cerberus-examples/internal/database"
	"cerberus-examples/internal/repositories"
	"context"
	"fmt"
)

type StoryService interface {
	Create(ctx context.Context, sprintId, description string) (repositories.Story, error)
	FindBySprint(ctx context.Context, sprintId string) ([]repositories.Story, error)
	Get(ctx context.Context, storyId string) (repositories.Story, error)
	Assign(ctx context.Context, storyId, userId string) (repositories.Story, error)
	Estimate(ctx context.Context, storyId string, estimation int) (repositories.Story, error)
	ChangeStatus(ctx context.Context, storyId, status string) (repositories.Story, error)
}

type storyService struct {
	txProvider database.TxProvider
	repo       repositories.StoryRepo
}

func NewStoryService(
	txProvider database.TxProvider,
	repo repositories.StoryRepo) StoryService {
	return &storyService{
		txProvider: txProvider,
		repo:       repo,
	}
}

func (s *storyService) Create(ctx context.Context, sprintId, description string) (repositories.Story, error) {

	userId := ctx.Value("userId")
	if userId == nil {
		return repositories.Story{}, fmt.Errorf("no userId")
	}

	tx, err := s.txProvider.GetTransaction()
	if err != nil {
		return repositories.Story{}, err
	}

	story, err := s.repo.Create(sprintId, description, tx)
	if err != nil {
		if rbe := tx.Rollback(); rbe != nil {
			err = fmt.Errorf("rollback error (%v) after %w", rbe, err)
		}
		return repositories.Story{}, err
	}

	return story, tx.Commit()
}

func (s *storyService) FindBySprint(ctx context.Context, sprintId string) ([]repositories.Story, error) {
	return s.repo.FindBySprint(sprintId)
}

func (s *storyService) Get(ctx context.Context, storyId string) (repositories.Story, error) {
	return s.repo.Get(storyId, nil)
}

func (s *storyService) Assign(ctx context.Context, storyId, userId string) (repositories.Story, error) {
	_, err := s.repo.Assign(storyId, userId)
	if err != nil {
		return repositories.Story{}, err
	}
	return s.repo.Get(storyId, nil)
}

func (s *storyService) Estimate(ctx context.Context, storyId string, estimation int) (repositories.Story, error) {
	_, err := s.repo.Estimate(storyId, estimation)
	if err != nil {
		return repositories.Story{}, err
	}
	return s.repo.Get(storyId, nil)
}

func (s *storyService) ChangeStatus(ctx context.Context, storyId, status string) (repositories.Story, error) {
	_, err := s.repo.ChangeStatus(storyId, status)
	if err != nil {
		return repositories.Story{}, err
	}
	return s.repo.Get(storyId, nil)
}
