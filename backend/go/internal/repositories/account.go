package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"log"
)

type AccountRepo interface {
	Create(tx *sql.Tx) (Account, error)
}

type Account struct {
	Id string `json:"id"`
}

type accountRepo struct {
	db *sql.DB
}

func NewAccountRepo(db *sql.DB) AccountRepo {
	return &accountRepo{
		db: db,
	}
}

func (r *accountRepo) Create(tx *sql.Tx) (account Account, err error) {
	if tx != nil {
		return r.create(tx)
	}

	tx, err = r.db.Begin()
	if err != nil {
		log.Println(err)
		return
	}

	account, err = r.create(tx)

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (r *accountRepo) create(tx *sql.Tx) (account Account, err error) {
	stmt, err := tx.Prepare("insert into account(id) values(?)")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	id := uuid.New().String()
	_, err = stmt.Exec(id)
	if err != nil {
		log.Println(err)
		return
	}

	account = Account{
		Id: id,
	}
	return
}
