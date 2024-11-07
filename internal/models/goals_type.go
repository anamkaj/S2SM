package models

type GoalsList interface {
	GetGoalsCounter(id int) (*TransformGoals, error)
	CheckCounter(id int) bool
	InsertGoal(data *TransformGoals, info *CounterInfo) error
	UpdateGoalStatus(goals *[]GoalsDataBase) error
	DeleteCounter(count int) error
}

type TransformGoals struct {
	CounterID int              `json:"counter_id"`
	Goals     *[]GoalsDataBase `json:"goals"`
}

type GoalsDataBase struct {
	GoalID int    `json:"goal_id" db:"goal_id" validate:"required"`
	Name   string `json:"name" db:"name" validate:"required,min=1,max=255"`
	Status bool   `json:"status" db:"status"`
}

type CounterInfo struct {
	ID         int    `json:"id" db:"counter_id"`
	Status     string `json:"status" db:"status"`
	OwnerLogin string `json:"owner_login" db:"owner_login"`
	Name       string `json:"name" db:"name"`
	Site       string `json:"site" db:"site"`
	Site2      Site2
}

type Site2 struct {
	Site   string `json:"site" db:"site_two"`
	Domain string `json:"domain" db:"domain"`
}
