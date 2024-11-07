package models

type ResStatRequest struct {
	Query struct {
		Ids        []int    `json:"ids"`
		Dimensions []string `json:"dimensions"`
		Metrics    []string `json:"metrics"`
		Date1      string   `json:"date1"`
		Date2      string   `json:"date2"`
	} `json:"query"`
	Data []struct {
		Dimensions []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"dimensions"`
		Metrics []float64 `json:"metrics"`
	} `json:"data"`
	Totals []float64 `json:"totals"`
}

type TransformStats struct {
	GoalID  int     `json:"goal_id" `
	Name    string  `json:"name" `
	Status  bool    `json:"status" `
	Metrics float64 `json:"metrics" `
}

type ResultTransformStats struct {
	Data1    string            `json:"data1"`
	Data2    string            `json:"data2"`
	LenGoals int               `json:"len_goals"`
	Goal     *[]TransformStats `json:"goal"`
}

type ResGoalRequest struct {
	Goals []struct {
		ID            int     `json:"id"`
		Name          string  `json:"name"`
		Type          string  `json:"type"`
		DefaultPrice  float64 `json:"default_price"`
		IsRetargeting int     `json:"is_retargeting"`
		GoalSource    string  `json:"goal_source"`
		IsFavorite    int     `json:"is_favorite"`
		PrevGoalID    int     `json:"prev_goal_id,omitempty"`
		Conditions    []struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"conditions,omitempty"`
		Steps []struct {
			ID            int     `json:"id"`
			Name          string  `json:"name"`
			Type          string  `json:"type"`
			DefaultPrice  float64 `json:"default_price"`
			IsRetargeting int     `json:"is_retargeting"`
			GoalSource    string  `json:"goal_source"`
			IsFavorite    int     `json:"is_favorite"`
			PrevGoalID    int     `json:"prev_goal_id"`
			Conditions    []struct {
				Type string `json:"type"`
				URL  string `json:"url"`
			} `json:"conditions"`
			PriceLocations []interface{} `json:"price_locations"`
		} `json:"steps,omitempty"`
		Flag            string        `json:"flag,omitempty"`
		PriceLocations  []interface{} `json:"price_locations,omitempty"`
		HidePhoneNumber bool          `json:"hide_phone_number,omitempty"`
	} `json:"goals"`
}
