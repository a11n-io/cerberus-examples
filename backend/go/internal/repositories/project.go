package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"log"
)

type ProjectRepo interface {
	Create(accountId, name, description string, tx *sql.Tx) (Project, error)
	FindByAccount(accountId string) ([]Project, error)
	Get(projectId string) (Project, error)
	Delete(projectId string) error
}

type Project struct {
	Id          string `json:"id"`
	AccountId   string `json:"accountId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type projectRepo struct {
	db *sql.DB
}

func NewProjectRepo(db *sql.DB) ProjectRepo {
	return &projectRepo{
		db: db,
	}
}

func (r *projectRepo) Create(accountId, name, description string, tx *sql.Tx) (project Project, err error) {

	if tx != nil {
		return r.create(accountId, name, description, tx)
	}

	tx, err = r.db.Begin()
	if err != nil {
		log.Println(err)
		return
	}

	project, err = r.create(accountId, name, description, tx)
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

func (r *projectRepo) create(accountId, name, description string, tx *sql.Tx) (project Project, err error) {
	stmt, err := tx.Prepare("insert into project(id, account_id, name, description) values(?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	id := uuid.New().String()
	_, err = stmt.Exec(id, accountId, name, description)
	if err != nil {
		log.Println(err)
		return
	}

	project = Project{
		Id:          id,
		AccountId:   accountId,
		Name:        name,
		Description: description,
	}
	return
}

func (r *projectRepo) FindByAccount(accountId string) (projects []Project, err error) {

	stmt, err := r.db.Prepare("select id, name, description from project where account_id = ? order by name asc")
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
		var id, name, description string
		err = rows.Scan(&id, &name, &description)
		if err != nil {
			return
		}

		projects = append(projects, Project{
			Id:          id,
			AccountId:   accountId,
			Name:        name,
			Description: description,
		})
	}

	return
}

func (r *projectRepo) Get(projectId string) (project Project, err error) {

	stmt, err := r.db.Prepare("select account_id, name, description from project where id = ?")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	var name, description, accountId string
	err = stmt.QueryRow(projectId).Scan(&accountId, &name, &description)
	if err != nil {
		return
	}

	project = Project{
		Id:          projectId,
		AccountId:   accountId,
		Name:        name,
		Description: description,
	}

	return
}

func (r *projectRepo) Delete(projectId string) (err error) {

	stmt, err := r.db.Prepare("delete from project where id = ?")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(projectId)

	return
}
