package ports

import (
	"context"

	"github.com/mayank18049/GoChessMoveShare/internal/domain/aggregate"
)

type GameRepoStatus int

const (
	GAMEREPO_OK GameRepoStatus = -iota
	GAMEREPO_FAILED
	GAMEREPO_ENTRY_NOT_FOUND
	GAMEREPO_ENTRY_EXISTS
	GAMEREPO_NOT_IMPLEMENTED
)

type GameRepo interface {
	CreateGame(ctx context.Context, game aggregate.Game) (string, GameRepoStatus)
	GetGame(ctx context.Context, gameID string) (aggregate.Game, GameRepoStatus)
	SetGame(ctx context.Context, gameID string, game aggregate.Game) GameRepoStatus
	DeleteGame(ctx context.Context, gameID string) GameRepoStatus
}
