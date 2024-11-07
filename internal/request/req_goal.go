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

type GoalsRequest interface {
	GetGoalRequest(id int, db *models.GoalsList) (*models.ResGoalRequest, error)
	InfoCounter(id int) (*models.CounterInfo, error)
	InsertCounter(id int) error
}

type Request struct {
	db models.GoalsList
}

func NewRequestGoal(db models.GoalsList) *Request {
	return &Request{
		db: db,
	}
}

func (r *Request) GetGoalRequest(id int) (*models.TransformGoals, error) {

	token, err := utils.GetToken()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}
	url := fmt.Sprintf("https:/%s/management/v1/counter/%d/goals", token.UrlApi, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err

	}
	req.Header.Set("Content-Type", "application/x-yametrika+json")
	req.Header.Set("Authorization", token.AccessToken)
	req.Header.Set("Cookie", token.CookieSession)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s ошибка запроса к серверу", err)
		return nil, err

	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err

	}

	var response models.ResGoalRequest

	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err

	}

	transform := utils.GoalsTransform(&response, id)

	return transform, err

}

func (r *Request) InfoCounter(id int) (*models.CounterInfo, error) {

	token, err := utils.GetToken()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}
	url := fmt.Sprintf("https://%s/management/v1/counter/%d", token.UrlApi, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error: %s", err)

	}
	req.Header.Set("Content-Type", "application/x-yametrika+json")
	req.Header.Set("Authorization", token.AccessToken)
	req.Header.Set("Cookie", token.CookieSession)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error: %s ошибка запроса к серверу - информация о счетчике", err)

	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error: %s ошибка запроса к серверу - информация о счетчике", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)

	}

	var counterData struct {
		Counter models.CounterInfo `json:"counter"`
	}

	if err := json.Unmarshal(body, &counterData); err != nil {
		fmt.Printf("Error:ошибка десериализации %s", err)
	}

	response := counterData.Counter

	return &response, err

}

func (r *Request) InsertCounter(id int) error {

	goals, err := r.GetGoalRequest(id)
	if err != nil {
		return err
	}

	info, err := r.InfoCounter(id)
	if err != nil {
		return err
	}

	err = r.db.InsertGoal(goals, info)
	if err != nil {
		return fmt.Errorf("error: ошибка добавления целей в базу %s", err)
	}

	return nil

}
