package repositories

import (
	"database/sql"
	"fmt"
	cerberus "github.com/a11n-io/go-cerberus"
	"github.com/google/uuid"
	"log"
)

type UserRepo interface {
	Save(accountId, email, plainPassword, name string, tx *sql.Tx) (User, error)
	FindOneByEmailAndPassword(email string, password string) (User, error)
	FindAll(accountId string) ([]User, error)
}

type User struct {
	Token             string             `json:"token"`
	CerberusTokenPair cerberus.TokenPair `json:"cerberusTokenPair"`
	Id                string             `json:"id"`
	AccountId         string             `json:"accountId"`
	Name              string             `json:"name"`
	Email             string             `json:"email"`
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Save(accountId, email, plainPassword, name string, tx *sql.Tx) (user User, err error) {

	encryptedPassword, err := encryptPassword(plainPassword)
	if err != nil {
		log.Println(err)
		return
	}

	if tx != nil {
		return r.save(accountId, email, encryptedPassword, name, tx)
	}

	tx, err = r.db.Begin()
	if err != nil {
		log.Println(err)
		return
	}

	user, err = r.save(accountId, email, encryptedPassword, name, tx)
	if err != nil {
		log.Println(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (r *userRepo) save(accountId, email, encryptedPassword, name string, tx *sql.Tx) (user User, err error) {

	stmt, err := tx.Prepare("insert into user(id, account_id, email, password, name) values(?, ?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	id := uuid.New().String()
	_, err = stmt.Exec(id, accountId, email, encryptedPassword, name)
	if err != nil {
		log.Println(err)
		return
	}

	user = User{
		Id:        id,
		AccountId: accountId,
		Name:      name,
		Email:     email,
	}
	return
}

func (r *userRepo) FindOneByEmailAndPassword(email string, plainPassword string) (user User, err error) {

	stmt, err := r.db.Prepare("select id, account_id, name, password from user where email = ?")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	var id, accountId, name, password string
	err = stmt.QueryRow(email).Scan(&id, &accountId, &name, &password)
	if err != nil {
		err = fmt.Errorf("account not found or incorrect password")
		return
	}

	if !verifyPassword(plainPassword, password) {
		err = fmt.Errorf("account not found or incorrect password")
		return
	}

	user = User{
		Id:        id,
		AccountId: accountId,
		Name:      name,
		Email:     email,
	}

	return
}

func (r *userRepo) FindAll(accountId string) (users []User, err error) {

	stmt, err := r.db.Prepare("select id, name, email from user where account_id = ? order by name asc")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(accountId)
	if err != nil {
		return
	}

	for rows.Next() {
		var id, name, email string
		err = rows.Scan(&id, &name, &email)
		if err != nil {
			return
		}

		users = append(users, User{
			Id:        id,
			AccountId: accountId,
			Name:      name,
			Email:     email,
		})
	}

	return
}
