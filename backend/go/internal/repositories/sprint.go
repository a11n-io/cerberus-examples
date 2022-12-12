package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"log"
	"time"
)

type SprintRepo interface {
	Create(projectId, goal string, tx *sql.Tx) (Sprint, error)
	FindByProject(projectId string) ([]Sprint, error)
	Get(sprintId string, tx *sql.Tx) (Sprint, error)
	Start(sprintId string) (Sprint, error)
	End(sprintId string) (Sprint, error)
}

type Sprint struct {
	Id           string `json:"id"`
	ProjectId    string `json:"projectId"`
	SprintNumber int    `json:"sprintNumber"`
	Goal         string `json:"goal"`
	StartDate    int64  `json:"startDate"`
	EndDate      int64  `json:"endDate"`
}

type sprintRepo struct {
	db *sql.DB
}

func NewSprintRepo(db *sql.DB) SprintRepo {
	return &sprintRepo{
		db: db,
	}
}

func (r *sprintRepo) Create(projectId, goal string, tx *sql.Tx) (sprint Sprint, err error) {
	if tx != nil {
		return r.create(projectId, goal, tx)
	}

	tx, err = r.db.Begin()
	if err != nil {
		log.Println(err)
		return
	}

	sprint, err = r.create(projectId, goal, tx)
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

func (r *sprintRepo) create(projectId, goal string, tx *sql.Tx) (sprint Sprint, err error) {
	stmt, err := tx.Prepare("insert into sprint(id, project_id, sprint_number, goal, start_date, end_date)" +
		" values(?, ?, " +
		"(SELECT COUNT(*) + 1 FROM sprint WHERE project_id = ?), " +
		"?, 0, 0)")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	id := uuid.New().String()
	_, err = stmt.Exec(id, projectId, projectId, goal)
	if err != nil {
		log.Println(err)
		return
	}

	return r.Get(id, tx)
}

func (r *sprintRepo) FindByProject(projectId string) (sprints []Sprint, err error) {

	stmt, err := r.db.Prepare(
		"select id, sprint_number, goal, start_date, end_date from sprint " +
			"where project_id = ? order by sprint_number asc")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(projectId)
	if err != nil {
		return
	}

	for rows.Next() {
		var id, goal string
		var sprintNumber int
		var startDate, endDate int64
		err = rows.Scan(&id, &sprintNumber, &goal, &startDate, &endDate)
		if err != nil {
			return
		}

		sprints = append(sprints, Sprint{
			Id:           id,
			ProjectId:    projectId,
			SprintNumber: sprintNumber,
			Goal:         goal,
			StartDate:    startDate,
			EndDate:      endDate,
		})
	}

	return
}

func (r *sprintRepo) Get(sprintId string, tx *sql.Tx) (sprint Sprint, err error) {

	if tx != nil {
		return r.get(sprintId, tx)
	}

	tx, err = r.db.Begin()
	if err != nil {
		log.Println(err)
		return
	}

	sprint, err = r.get(sprintId, tx)
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

func (r *sprintRepo) get(sprintId string, tx *sql.Tx) (sprint Sprint, err error) {
	stmt, err := tx.Prepare("select project_id, sprint_number, goal, start_date, end_date from sprint where id = ?")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	var goal, projectId string
	var sprintNumber int
	var startDate, endDate int64
	err = stmt.QueryRow(sprintId).Scan(&projectId, &sprintNumber, &goal, &startDate, &endDate)
	if err != nil {
		return
	}

	sprint = Sprint{
		Id:           sprintId,
		ProjectId:    projectId,
		SprintNumber: sprintNumber,
		Goal:         goal,
		StartDate:    startDate,
		EndDate:      endDate,
	}
	return
}

func (r *sprintRepo) Start(sprintId string) (sprint Sprint, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	stmt, err := tx.Prepare("update sprint set start_date = ? where id = ?")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	startDate := time.Now().Unix()
	_, err = stmt.Exec(startDate, sprintId)
	if err != nil {
		log.Println(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}

	sprint = Sprint{
		Id:        sprintId,
		StartDate: startDate,
	}

	return
}

func (r *sprintRepo) End(sprintId string) (sprint Sprint, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	stmt, err := tx.Prepare("update sprint set end_date = ? where id = ?")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	endDate := time.Now().Unix()
	_, err = stmt.Exec(endDate, sprintId)
	if err != nil {
		log.Println(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}

	sprint = Sprint{
		Id:      sprintId,
		EndDate: endDate,
	}

	return
}
