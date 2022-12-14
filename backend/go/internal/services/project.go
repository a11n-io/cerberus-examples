package services

import (
	"cerberus-examples/internal/common"
	"cerberus-examples/internal/database"
	"cerberus-examples/internal/repositories"
	"context"
	"fmt"
	cerberus "github.com/a11n-io/go-cerberus"
)

type ProjectService interface {
	Create(ctx context.Context, accountId, name, description string) (repositories.Project, error)
	FindAll(ctx context.Context, accountId string) ([]repositories.Project, error)
	Get(ctx context.Context, projectId string) (repositories.Project, error)
	Delete(ctx context.Context, projectId string) error
}

type projectService struct {
	txProvider     database.TxProvider
	repo           repositories.ProjectRepo
	cerberusClient cerberus.CerberusClient
}

func NewProjectService(
	txProvider database.TxProvider,
	repo repositories.ProjectRepo,
	cerberusClient cerberus.CerberusClient) ProjectService {
	return &projectService{
		txProvider:     txProvider,
		repo:           repo,
		cerberusClient: cerberusClient,
	}
}

func (s *projectService) Create(ctx context.Context, accountId, name, description string) (repositories.Project, error) {

	userId := ctx.Value("userId")
	if userId == nil {
		return repositories.Project{}, fmt.Errorf("no userId")
	}

	tx, err := s.txProvider.GetTransaction()
	if err != nil {
		return repositories.Project{}, err
	}

	project, err := s.repo.Create(accountId, name, description, tx)
	if err != nil {
		if rbe := tx.Rollback(); rbe != nil {
			err = fmt.Errorf("rollback error (%v) after %w", rbe, err)
		}
		return repositories.Project{}, err
	}

	err = s.cerberusClient.ExecuteWithCtx(ctx, s.cerberusClient.CreateResourceCmd(project.Id, accountId, common.Project_RT))
	if err != nil {
		if rbe := tx.Rollback(); rbe != nil {
			err = fmt.Errorf("rollback error (%v) after %w", rbe, err)
		}
		return repositories.Project{}, err
	}

	return project, tx.Commit()
}

func (s *projectService) FindAll(ctx context.Context, accountId string) ([]repositories.Project, error) {
	return s.repo.FindByAccount(accountId)
}

func (s *projectService) Get(ctx context.Context, projectId string) (repositories.Project, error) {
	return s.repo.Get(projectId)
}

func (s *projectService) Delete(ctx context.Context, projectId string) error {
	return s.repo.Delete(projectId)
}
