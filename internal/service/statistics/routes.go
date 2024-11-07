package statistics

import (
	"encoding/json"
	"fmt"
	"io"
	"metrika/internal/models"
	"metrika/internal/request"
	"metrika/internal/utils"
	"net/http"
	"github.com/go-playground/validator/v10"
)

type ParamRequest struct {
	Id         int    `json:"id_count" validate:"required"`
	Date_start string `json:"date_start" validate:"required,min=10,max=10" format:"2006-01-02" `
	Date_end   string `json:"date_end" validate:"required,min=10,max=10" format:"2006-01-02"`
}

type Handler struct {
	store models.GoalsList
}

func NewHandler(store models.GoalsList) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.Handle("POST /api/metrika/statistics", http.HandlerFunc(h.handleGetStatistics))
}

func (h *Handler) handleGetStatistics(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var statistics *ParamRequest

	err = json.Unmarshal(body, &statistics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(statistics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	params := &request.Param{
		Id:    statistics.Id,
		Start: statistics.Date_start,
		End:   statistics.Date_end,
	}

	request := request.NewRequestStat(h.store)
	goals, err := request.GetStatApi(params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	utils.ResJson(w, http.StatusOK, goals)

}
