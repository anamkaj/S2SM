package api

import (
	"fmt"
	"log"
	"metrika/internal/config"
	"metrika/internal/service/goals"
	"metrika/internal/service/statistics"
	"metrika/pkg/logger"
	"metrika/pkg/logger/middleware"
	"net/http"
	"github.com/jmoiron/sqlx"
)

type ApiServer struct {
	addr string
	db   *sqlx.DB
}

func NewApiServer(addr string, db *sqlx.DB) *ApiServer {
	return &ApiServer{
		addr: addr,
		db:   db,
	}
}

func (s *ApiServer) Run() error {
	mux := http.NewServeMux()

	storeGoals := goals.NewStore(s.db)
	goalsHandler := goals.NewHandler(storeGoals)
	goalsHandler.RegisterRoutes(mux)

	statHandler := statistics.NewHandler(storeGoals)
	statHandler.RegisterRoutes(mux)

	loggerJournal, err := logger.NewLogger(s.db)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	LoggerMux := middleware.LoggerMiddleware(mux, loggerJournal)
	cors := config.Cors(LoggerMux)

	if err := loggerJournal.LoggerBasic(logger.INFO_LOG, "Server started on port 8060"); err != nil {
		log.Println("Failed to log to database:", err)
	}
	fmt.Println("Server started on port 8080")

	if err := http.ListenAndServe(":8080", cors); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return fmt.Errorf("server error")
}
