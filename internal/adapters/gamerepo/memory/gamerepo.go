package memory

import (
	"context"
	"log"

	"github.com/mayank18049/GoChessMoveShare/internal/domain/aggregate"
	"github.com/mayank18049/GoChessMoveShare/internal/domain/ports"
)

type InMemoryGameRepo struct {
	games           map[string]aggregate.Game
	currentTeachers map[string]string
	logger          *log.Logger
}

func NewInMemoryGameRepo(logger *log.Logger) *InMemoryGameRepo {
	return &InMemoryGameRepo{
		games:           make(map[string]aggregate.Game),
		currentTeachers: make(map[string]string),
		logger:          logger,
	}
}

func (gr *InMemoryGameRepo) CreateGame(ctx context.Context, game aggregate.Game) (string, ports.GameRepoStatus) {
	gr.logger.Printf("[MemoryGameRepo CreateGame]: Creating new game in repo with Game ID %s and teacher ID %s/n", game.GetID(), game.GetTeacherID())
	tid := game.GetTeacherID()
	gameID, present := gr.currentTeachers[tid]
	if present {
		return gameID, ports.GAMEREPO_ENTRY_EXISTS
	}
	gameID = game.GetID()
	gr.currentTeachers[tid] = gameID
	gr.games[gameID] = game
	gr.logger.Printf("[MemoryGameRepo CreateGame]: Printing updated and teacher ID map %s/n", gr.currentTeachers)
	return gameID, ports.GAMEREPO_OK
}

func (gr *InMemoryGameRepo) DeleteGame(ctx context.Context, gameID string) ports.GameRepoStatus {
	gr.logger.Printf("[MemoryGameRepo DeleteGame]: Deleting ID %s\n", gameID)
	game, present := gr.games[gameID]

	if !present {
		gr.logger.Printf("[MemoryGameRepo DeleteGame]: Delete Failed %s\n", gr.currentTeachers)
		return ports.GAMEREPO_ENTRY_NOT_FOUND
	}
	tid := game.GetTeacherID()
	delete(gr.currentTeachers, tid)
	delete(gr.games, gameID)
	return ports.GAMEREPO_OK
}

func (gr *InMemoryGameRepo) GetGame(ctx context.Context, gameID string) (aggregate.Game, ports.GameRepoStatus) {
	gr.logger.Printf("[MemoryGameRepo GetGame]: Getting ID %s\n", gameID)
	game, present := gr.games[gameID]
	if !present {
		gr.logger.Printf("[MemoryGameRepo GetGame]: Failed Getting ID %s\n", gr.currentTeachers)
		return aggregate.Game{}, ports.GAMEREPO_ENTRY_NOT_FOUND
	}
	return game, ports.GAMEREPO_OK
}

func (gr *InMemoryGameRepo) SetGame(ctx context.Context, gameID string, game aggregate.Game) ports.GameRepoStatus {
	gr.logger.Printf("[MemoryGameRepo SetGame]: Setting game for ID %s\n", gameID)
	_, present := gr.games[gameID]
	if !present {
		gr.logger.Printf("[MemoryGameRepo SetGame]: GameID not found %s\n", gameID)
		return ports.GAMEREPO_ENTRY_NOT_FOUND
	}
	gr.games[gameID] = game
	return ports.GAMEREPO_OK
}
