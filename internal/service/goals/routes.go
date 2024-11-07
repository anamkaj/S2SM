package goals

import (
	"encoding/json"
	"fmt"
	"io"
	"metrika/internal/models"
	"metrika/internal/request"
	"metrika/internal/utils"
	"net/http"
	"strconv"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store models.GoalsList
}

func NewHandler(store models.GoalsList) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.Handle("GET /api/metrika/goals", http.HandlerFunc(h.handleGetGoals))
	router.Handle("PUT /api/metrika/update", http.HandlerFunc(h.handleUpdateStatusGoal))
	router.Handle("DELETE /api/metrika/delete-count", http.HandlerFunc(h.handleDeleteCounter))
}

func (h *Handler) handleGetGoals(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id_count")
	if id == "" || len(id) <= 3 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Ошибка: отсутствует параметр id_count или слишком короткий")
		return
	}
	num, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status := h.store.CheckCounter(num)
	if !status {
		err := request.NewRequestGoal(h.store).InsertCounter(num)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}
	}

	goals, err := h.store.GetGoalsCounter(num)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	utils.ResJson(w, http.StatusOK, goals)
	return

}

func (h *Handler) handleUpdateStatusGoal(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var goal struct {
		Goals []models.GoalsDataBase `json:"goals"`
	}

	err = json.Unmarshal(body, &goal)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(goal.Goals) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", "Отсутствуют цели")
		return
	}

	validate := validator.New()
	for _, v := range goal.Goals {
		err = validate.Struct(v)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}
	}

	err = h.store.UpdateGoalStatus(&goal.Goals)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}
	utils.ResJson(w, http.StatusOK, "Goals updated")

}

func (h *Handler) handleDeleteCounter(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id_count")
	if id == "" || len(id) <= 3 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Ошибка: отсутствует параметр id_count или слишком короткий")
		return
	}

	num, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status := h.store.CheckCounter(num)
	if !status {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Ошибка: счетчик не существует")
		return
	}

	err = h.store.DeleteCounter(num)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}
	utils.ResJson(w, http.StatusOK, fmt.Sprintf("Counter %d deleted", num))
}
