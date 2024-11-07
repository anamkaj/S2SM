package goals

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"metrika/internal/models"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetGoalsCounter(id int) (*models.TransformGoals, error) {

	qGoals := `WITH counter_check AS (
            SELECT goal_id, name, status
            FROM goals
            WHERE fk_counters_metrika_counter_id = $1)
            SELECT goal_id, name, status FROM counter_check;`

	data := []models.GoalsDataBase{}

	err := s.db.Select(&data, qGoals, id)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("на %d счетчике нет целей", id)
	}

	result := models.TransformGoals{
		CounterID: id,
		Goals:     &data,
	}

	return &result, err

}

func (s *Store) InsertGoal(data *models.TransformGoals, info *models.CounterInfo) error {

	qInsertCounter := `INSERT INTO counters_metrika (counter_id, status, owner_login, name, site, site_two, domain)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (counter_id)
        DO UPDATE SET
            status = EXCLUDED.status,
            owner_login = EXCLUDED.owner_login,
            name = EXCLUDED.name,
            site = EXCLUDED.site,
            site_two = EXCLUDED.site_two,
            domain = EXCLUDED.domain;`

	qInsertGoals := `INSERT INTO goals (goal_id, name, fk_counters_metrika_counter_id,status)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (goal_id)
            DO UPDATE SET
                goal_id = EXCLUDED.goal_id,
                name = EXCLUDED.name,
                fk_counters_metrika_counter_id = EXCLUDED.fk_counters_metrika_counter_id,
                status = EXCLUDED.status;`

	_, err := s.db.Exec(qInsertCounter, info.ID, info.Status, info.OwnerLogin, info.Name, info.Site, info.Site2.Site, info.Site2.Domain)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}

	for _, goal := range *data.Goals {
		_, err := s.db.Exec(qInsertGoals, goal.GoalID, goal.Name, data.CounterID, goal.Status)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return err
		}
	}

	return err

}

func (s *Store) UpdateGoalStatus(goals *[]models.GoalsDataBase) error {

	qUpdateGoals := `UPDATE goals
        SET
        status = $2
        WHERE goal_id = $1;`

	for _, goal := range *goals {
		_, err := s.db.Exec(qUpdateGoals, goal.GoalID, goal.Status)
		if err != nil {
			return fmt.Errorf("error: %s", err)
		}
	}
	return nil

}

func (s *Store) DeleteCounter(count int) error {
	qDeleteCounter := `WITH deleted_counters AS (
        DELETE FROM counters_metrika
        WHERE counter_id = $1
        RETURNING counter_id
    )
    DELETE FROM goals
    WHERE fk_counters_metrika_counter_id IN (SELECT counter_id FROM deleted_counters);`

	_, err := s.db.Exec(qDeleteCounter, count)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	return nil
}

func (s *Store) CheckCounter(id int) bool {
	qCheckCounter := `SELECT counter_id FROM counters_metrika WHERE counter_id = $1;`

	err := s.db.Get(&id, qCheckCounter, id)
	if err != nil {
		return false
	} else {
		return false
	}

}
