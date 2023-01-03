package services

import (
	"cerberus-examples/internal/database"
	"cerberus-examples/internal/repositories"
	"cerberus-examples/internal/services/jwtutils"
	"context"
	"fmt"
	"log"
)

type UserService interface {
	Register(ctx context.Context, email, plainPassword, name string) (repositories.User, error)
	Login(ctx context.Context, email string, password string) (repositories.User, error)
	Add(ctx context.Context, email, plainPassword, name, roleId string) (repositories.User, error)
	GetAll(ctx context.Context) ([]repositories.User, error)
}

type userService struct {
	txProvider  database.TxProvider
	userRepo    repositories.UserRepo
	accountRepo repositories.AccountRepo
	jwtSecret   string
	saltRounds  int
}

func NewUserService(
	txProvider database.TxProvider,
	userRepo repositories.UserRepo,
	accountRepo repositories.AccountRepo,
	jwtSecret string,
	saltRounds int) UserService {
	return &userService{
		txProvider:  txProvider,
		userRepo:    userRepo,
		accountRepo: accountRepo,
		jwtSecret:   jwtSecret,
		saltRounds:  saltRounds,
	}
}

// Register should register a new user
//
// The properties also be used to generate a JWT `token` which should be included
// with the returned user.
func (s *userService) Register(ctx context.Context, email, plainPassword, name string) (_ repositories.User, err error) {

	log.Println("Register", email, plainPassword, name)

	tx, err := s.txProvider.GetTransaction()
	if err != nil {
		return repositories.User{}, err
	}

	account, err := s.accountRepo.Create(tx)
	if err != nil {
		if rbe := tx.Rollback(); rbe != nil {
			err = fmt.Errorf("rollback error (%v) after %w", rbe, err)
		}
		return repositories.User{}, err
	}

	user, err := s.userRepo.Save(account.Id, email, plainPassword, name, tx)
	if err != nil {
		if rbe := tx.Rollback(); rbe != nil {
			err = fmt.Errorf("rollback error (%v) after %w", rbe, err)
		}
		return repositories.User{}, err
	}

	subject := user.Id
	token, err := jwtutils.Sign(subject, toClaims(user), s.jwtSecret)
	if err != nil {
		if rbe := tx.Rollback(); rbe != nil {
			err = fmt.Errorf("rollback error (%v) after %w", rbe, err)
		}
		return repositories.User{}, err
	}

	return userWithTokens(user, token), tx.Commit()
}

// Login finds a user and returns that user with a jwt token
func (s *userService) Login(ctx context.Context, email string, password string) (_ repositories.User, err error) {

	user, err := s.userRepo.FindOneByEmailAndPassword(email, password)
	if err != nil {
		return repositories.User{}, err
	}

	subject := user.Id
	token, err := jwtutils.Sign(subject, toClaims(user), s.jwtSecret)
	if err != nil {
		return repositories.User{}, err
	}

	return userWithTokens(user, token), nil
}

func (s *userService) Add(ctx context.Context, email, plainPassword, name, roleId string) (_ repositories.User, err error) {

	accountId := ctx.Value("accountId")
	if accountId == nil {
		return repositories.User{}, fmt.Errorf("no accountId")
	}

	tx, err := s.txProvider.GetTransaction()
	if err != nil {
		return repositories.User{}, err
	}

	user, err := s.userRepo.Save(accountId.(string), email, plainPassword, name, tx)
	if err != nil {
		if rbe := tx.Rollback(); rbe != nil {
			err = fmt.Errorf("rollback error (%v) after %w", rbe, err)
		}
		return repositories.User{}, err
	}

	return user, tx.Commit()
}

func (s *userService) GetAll(ctx context.Context) (_ []repositories.User, err error) {

	accountId := ctx.Value("accountId")
	if accountId == nil {
		return []repositories.User{}, fmt.Errorf("no accountId")
	}

	return s.userRepo.FindAll(accountId.(string))
}

func toClaims(user repositories.User) map[string]interface{} {
	return map[string]interface{}{
		"sub":       user.Id,
		"email":     user.Email,
		"name":      user.Name,
		"accountId": user.AccountId,
	}
}

func userWithTokens(user repositories.User, token string) repositories.User {
	return repositories.User{
		Token:     token,
		Id:        user.Id,
		AccountId: user.AccountId,
		Email:     user.Email,
		Name:      user.Name,
	}
}
