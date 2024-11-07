package utils

import (
	"fmt"
	"metrika/internal/models"
	"net/url"
	"strings"
	"time"
)

func FormatDate(start, end string) (string, string, error) {
	const layout = "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		return "", "", fmt.Errorf("ошибка парсинга начальной даты: %w", err)
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		return "", "", fmt.Errorf("ошибка парсинга конечной даты: %w", err)
	}
	return startDate.Format(layout), endDate.Format(layout), nil
}

func BuildUrl(id int, start string, end string, goals []models.GoalsDataBase) (string, error) {
	startDate, endDate, err := FormatDate(start, end)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return "", err
	}

	dimensions := "ym:s:goalDimension"
	filters := url.QueryEscape("(ym:s:AUTOMATICAdvEngine=='ya_direct' or ym:s:AUTOMATICAdvEngine=='ya_undefined') and (ym:s:isRobot=='No')")
	urlReq := fmt.Sprintf("https://api-metrika.yandex.ru/stat/v1/data?date1=%s&date2=%s&ids=%d&filters=%s&accuracy=full&dimensions=%s&metrics=", startDate, endDate, id, filters, dimensions)

	metrics := make([]string, len(goals))
	for index, goal := range goals {
		metrics[index] = fmt.Sprintf("ym:s:goal%dvisits", goal.GoalID)
	}
	urlReq += strings.Join(metrics, ",")
	return urlReq, err
}
