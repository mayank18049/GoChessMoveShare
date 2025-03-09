package handlers

import (
	"log"
	"net/http"

	"github.com/mayank18049/GoChessMoveShare/internal/DTO"
	"github.com/mayank18049/GoChessMoveShare/internal/service"
)

type GameHandler struct {
	regService service.GameRegistration
	logger     *log.Logger
}

func NewGameHandler(regService service.GameRegistration, logger *log.Logger) *GameHandler {
	return &GameHandler{
		regService: regService,
		logger:     logger}
}
func (gh *GameHandler) CreateGame(rw http.ResponseWriter, r *http.Request) {
	var gamereq DTO.CreateGameRequest
	err := FromJSON(&gamereq, r.Body)
	gh.logger.Printf("[GameHandler CreateGame]: %s\n", gamereq.TeacherID)
	if err != nil {
		gh.logger.Printf("[CreateGame]: Failed to decode JSON struct, %s\n", err.Error())
		http.Error(rw, "Failed to decode JSON struct", http.StatusBadRequest)
		return
	}
	gameresp, service_err := gh.regService.CreateGame(r.Context(), gamereq)
	if service_err != service.GAME_SERVICE_OK {
		gh.logger.Printf("[CreateGame]: Failed to Create Game, %d\n", service_err)
		http.Error(rw, "Failed to Create Game", http.StatusInternalServerError)
		return
	}
	ToJSON(gameresp, rw)

}

func (gh *GameHandler) DeleteGame(rw http.ResponseWriter, r *http.Request) {
	var gamereq DTO.DeleteGameRequest

	err := FromJSON(&gamereq, r.Body)
	gh.logger.Printf("[GameHandler DeleteGame]: %s\n", gamereq.GameID)
	if err != nil {
		gh.logger.Printf("[DeleteGame]: Failed to decode JSON struct, %s\n", err.Error())
		http.Error(rw, "Failed to decode JSON struct", http.StatusBadRequest)
		return
	}
	service_err := gh.regService.DeleteGame(r.Context(), gamereq)
	if service_err != service.GAME_SERVICE_OK {
		gh.logger.Printf("[DeleteGame]: Failed to Delete Game %d\n", service_err)
		http.Error(rw, "Failed to Delete Game", http.StatusInternalServerError)
		return
	}
	ToJSON("{}", rw)
}

func (gh *GameHandler) ConnectGame(rw http.ResponseWriter, r *http.Request) {
	var gamereq DTO.ConnectGameRequest
	err := FromJSON(&gamereq, r.Body)
	if err != nil {
		gh.logger.Printf("Failed to decode JSON struct, %s\n", err.Error())
		http.Error(rw, "Failed to decode JSON struct", http.StatusBadRequest)
		return
	}
	gameresp, service_err := gh.regService.ConnectGame(r.Context(), gamereq)
	if service_err != service.GAME_SERVICE_OK {
		gh.logger.Printf("Failed to Connect Game %d\n", service_err)
		http.Error(rw, "Failed to Connect Game", http.StatusInternalServerError)
		return
	}
	ToJSON(gameresp, rw)
}
