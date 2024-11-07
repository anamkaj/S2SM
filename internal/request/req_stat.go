package request

import (
	"encoding/json"
	"fmt"
	"io"
	"metrika/internal/models"
	"metrika/internal/utils"
	"net/http"
	"time"
)

type RequestStatApi interface {
	GetStatApi(params *Param) (*models.ResultTransformStats, error)
	TransformStats(params *TransformParams) (*[]models.TransformStats, error)
}

type Statistics struct {
	db models.GoalsList
}

type Param struct {
	Id    int
	Start string
	End   string
}

type TransformParams struct {
	Id    int
	Data  *models.ResStatRequest
	Goals []models.GoalsDataBase
}

func NewRequestStat(db models.GoalsList) *Statistics {
	return &Statistics{
		db: db,
	}
}

func (s *Statistics) GetStatApi(params *Param) (*models.ResultTransformStats, error) {

	token, err := utils.GetToken()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}

	status := s.db.CheckCounter(params.Id)
	if !status {
		return nil, fmt.Errorf("error: %d counter not found", params.Id)
	}

	goals, err := s.db.GetGoalsCounter(params.Id)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err.Error())
	}

	var statusActiveGoal []models.GoalsDataBase

	for _, goal := range *goals.Goals {
		if goal.Status {
			statusActiveGoal = append(statusActiveGoal, goal)
		}
	}

	if len(statusActiveGoal) == 0 {
		return nil, fmt.Errorf("на %d счетчике нет активных целей", params.Id)
	}

	// ________________________________________________________

	url, err := utils.BuildUrl(params.Id, params.Start, params.End, statusActiveGoal)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}

	// ________________________________________________________

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-yametrika+json")
	req.Header.Set("Authorization", token.AccessToken)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s ошибка запроса к серверу", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка: получен статистики %s", resp.Status)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}

	result := &models.ResStatRequest{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}

	// ________________________________________________________

	newParams := &TransformParams{
		Id:    params.Id,
		Data:  result,
		Goals: statusActiveGoal,
	}

	stat, err := s.TransformStats(newParams)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	data := models.ResultTransformStats{
		Data1:    result.Query.Date1,
		Data2:    result.Query.Date2,
		Goal:     stat,
		LenGoals: len(statusActiveGoal),
	}

	return &data, nil

}

func (s *Statistics) TransformStats(params *TransformParams) (*[]models.TransformStats, error) {

	var result []models.TransformStats

	for index, g := range params.Goals {

		if len(params.Data.Data) == 0 {
			return nil, fmt.Errorf("вернулись некорректные данные , массив пустой,попробуйте изменить диапазон дат")
		}

		total := params.Data.Totals[index]

		result = append(result, models.TransformStats{
			GoalID:  g.GoalID,
			Name:    g.Name,
			Status:  g.Status,
			Metrics: total,
		})

	}

	return &result, nil
}
