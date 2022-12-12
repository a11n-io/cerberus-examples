package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"log"
)

type StoryRepo interface {
	Create(sprintId, description string, tx *sql.Tx) (Story, error)
	FindBySprint(sprintId string) ([]Story, error)
	Get(storyId string, tx *sql.Tx) (Story, error)
	Estimate(storyId string, estimate int) (Story, error)
	ChangeStatus(storyId, status string) (Story, error)
	Assign(storyId, userId string) (Story, error)
}

type Story struct {
	Id          string `json:"id"`
	SprintId    string `json:"projectId"`
	Estimation  int    `json:"estimation"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Assignee    string `json:"assignee"`
}

type storyRepo struct {
	db *sql.DB
}

func NewStoryRepo(db *sql.DB) StoryRepo {
	return &storyRepo{
		db: db,
	}
}

func (r *storyRepo) Create(sprintId, description string, tx *sql.Tx) (story Story, err error) {

	if tx != nil {
		return r.create(sprintId, description, tx)
	}

	tx, err = r.db.Begin()
	if err != nil {
		log.Println(err)
		return
	}

	story, err = r.create(sprintId, description, tx)
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

func (r *storyRepo) create(sprintId, description string, tx *sql.Tx) (story Story, err error) {
	stmt, err := tx.Prepare("insert into story(id, sprint_id, estimation, description, status)" +
		" values(?, ?, 0, ?, 'todo')")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	id := uuid.New().String()
	_, err = stmt.Exec(id, sprintId, description)
	if err != nil {
		log.Println(err)
		return
	}

	story = Story{
		Id:          id,
		SprintId:    sprintId,
		Description: description,
		Estimation:  0,
		Status:      "todo",
	}

	return
}

func (r *storyRepo) FindBySprint(sprintId string) (stories []Story, err error) {

	stmt, err := r.db.Prepare(
		"select id, estimation, description, status, user_id from story " +
			"where sprint_id = ? order by description asc")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(sprintId)
	if err != nil {
		return
	}

	for rows.Next() {
		var id, description, status string
		var userId sql.NullString
		var estimation int
		err = rows.Scan(&id, &estimation, &description, &status, &userId)
		if err != nil {
			return
		}

		stories = append(stories, Story{
			Id:          id,
			SprintId:    sprintId,
			Estimation:  estimation,
			Description: description,
			Status:      status,
			Assignee:    userId.String,
		})
	}

	return
}

func (r *storyRepo) Get(storyId string, tx *sql.Tx) (story Story, err error) {
	if tx != nil {
		return r.get(storyId, tx)
	}

	tx, err = r.db.Begin()
	if err != nil {
		log.Println(err)
		return
	}

	story, err = r.get(storyId, tx)
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

func (r *storyRepo) get(storyId string, tx *sql.Tx) (story Story, err error) {
	stmt, err := tx.Prepare("select sprint_id, estimation, description, status, user_id from story where id = ?")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	var sprintId, description, status string
	var userId sql.NullString
	var estimation int
	err = stmt.QueryRow(storyId).Scan(&sprintId, &estimation, &description, &status, &userId)
	if err != nil {
		return
	}

	story = Story{
		Id:          storyId,
		SprintId:    sprintId,
		Estimation:  estimation,
		Description: description,
		Status:      status,
		Assignee:    userId.String,
	}

	return
}

func (r *storyRepo) Estimate(storyId string, estimation int) (story Story, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	stmt, err := tx.Prepare("update story set estimation = ? where id = ?")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(estimation, storyId)
	if err != nil {
		log.Println(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}

	story = Story{
		Id:         storyId,
		Estimation: estimation,
	}

	return
}

func (r *storyRepo) ChangeStatus(storyId, status string) (story Story, err error) {
	log.Println("ChangeStatus", storyId, status)
	tx, err := r.db.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	stmt, err := tx.Prepare("update story set status = ? where id = ?")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(status, storyId)
	if err != nil {
		log.Println(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}

	story = Story{
		Id:     storyId,
		Status: status,
	}

	return
}

func (r *storyRepo) Assign(storyId, userId string) (story Story, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	stmt, err := tx.Prepare("update story set user_id = ? where id = ?")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(userId, storyId)
	if err != nil {
		log.Println(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}

	story = Story{
		Id:       storyId,
		Assignee: userId,
	}

	return
}
