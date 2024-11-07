package utils

import (
	"fmt"
	"metrika/internal/models"
)

func GoalsTransform(goals *models.ResGoalRequest, id int) *models.TransformGoals {
	var data []models.GoalsDataBase

	for _, goal := range goals.Goals {
		if goal.Type == "action" || goal.Type == "form" ||
			goal.Type == "messenger" || goal.Type == "phone" ||
			goal.Type == "social" || goal.Type == "call" {

			data = append(data, models.GoalsDataBase{
				GoalID: goal.ID,
				Name:   goal.Name,
				Status: false,
			})
		}

		if goal.Type == "step" {
			name := goal.Name
			id := goal.ID
			data = append(data, models.GoalsDataBase{
				GoalID: id,
				Name:   name,
				Status: false,
			})

			for _, step := range goal.Steps {
				stepName := step.Name
				stepID := step.ID
				data = append(data, models.GoalsDataBase{
					GoalID: stepID,
					Name:   fmt.Sprintf("%s -> %s", name, stepName),
					Status: false,
				})
			}
		}
	}

	result := models.TransformGoals{
		CounterID: id,
		Goals:     &data,
	}

	return &result

}
