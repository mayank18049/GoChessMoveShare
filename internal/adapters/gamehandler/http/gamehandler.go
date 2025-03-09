package http

import (
	"log"
	"net/http"
	"os"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mayank18049/GoChessMoveShare/internal/adapters/gamehandler/http/handlers"
	"github.com/mayank18049/GoChessMoveShare/internal/service"
)

type HttpHandler struct {
	gameHandler handlers.GameHandler
	logger      *log.Logger
}

func NewHttpHandler(regService service.GameRegistration, logger *log.Logger) (*HttpHandler, error) {
	gameHandler := handlers.NewGameHandler(regService, logger)
	return &HttpHandler{
		gameHandler: *gameHandler,
		logger:      logger,
	}, nil
}

func (h *HttpHandler) SetAndServe(serverAddr string) {
	sm := mux.NewRouter()
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/create", h.gameHandler.CreateGame)
	postRouter.HandleFunc("/delete", h.gameHandler.DeleteGame)
	postRouter.HandleFunc("/connect", h.gameHandler.ConnectGame)
	cors := gorillaHandlers.CORS()
	s := http.Server{
		Addr:         serverAddr,
		Handler:      cors(sm),
		ErrorLog:     h.logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	err := s.ListenAndServe()
	if err != nil {
		h.logger.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
