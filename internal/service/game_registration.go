package service

import (
	"context"
	"log"

	"github.com/mayank18049/GoChessMoveShare/internal/DTO"
	"github.com/mayank18049/GoChessMoveShare/internal/domain/aggregate"
	"github.com/mayank18049/GoChessMoveShare/internal/domain/ports"
)

type GameServiceStatus int

const (
	GAME_SERVICE_OK GameServiceStatus = -iota
	GAME_SERVICE_FAILED
)

type GameRegistration struct {
	broker   ports.MessageBroker
	gamerepo ports.GameRepo
	logger   *log.Logger
}

func NewGameRegistration(broker ports.MessageBroker, gamerepo ports.GameRepo, logger *log.Logger) (*GameRegistration, GameServiceStatus) {
	return &GameRegistration{
		broker:   broker,
		gamerepo: gamerepo,
		logger:   logger}, GAME_SERVICE_OK
}

func (gr *GameRegistration) constructComms(ctx context.Context, game *aggregate.Game) GameServiceStatus {
	allocatedQueues, allocatedExchanges, allocatedConnections := make([]string, 0, 2), make([]string, 0, 2), make([][]string, 0, 1)

	status := gr.broker.CreateQueue(ctx, game.GetMoveQueueName(), ports.QUEUE_REPLAY)
	if status != ports.MESSAGE_SERVICE_OK {
		gr.destructComms(ctx, allocatedQueues, allocatedExchanges, allocatedConnections)
		return GAME_SERVICE_FAILED
	}
	allocatedQueues = append(allocatedQueues, game.GetMoveQueueName())
	status = gr.broker.CreateQueue(ctx, game.GetReplyQueueName(), ports.QUEUE_NO_REPLAY)
	if status != ports.MESSAGE_SERVICE_OK {
		gr.destructComms(ctx, allocatedQueues, allocatedExchanges, allocatedConnections)
		return GAME_SERVICE_FAILED
	}
	allocatedQueues = append(allocatedQueues, game.GetReplyQueueName())
	status = gr.broker.CreateExchange(ctx, game.GetControlExchangeName(), ports.EXCHANGE_FANOUT)
	if status != ports.MESSAGE_SERVICE_OK {
		gr.destructComms(ctx, allocatedQueues, allocatedExchanges, allocatedConnections)
		return GAME_SERVICE_FAILED
	}
	allocatedExchanges = append(allocatedExchanges, game.GetControlExchangeName())
	status = gr.broker.CreateExchange(ctx, game.GetReplyExchangeName(), ports.EXCHANGE_DIRECT)
	if status != ports.MESSAGE_SERVICE_OK {
		gr.destructComms(ctx, allocatedQueues, allocatedExchanges, allocatedConnections)
		return GAME_SERVICE_FAILED
	}
	allocatedExchanges = append(allocatedExchanges, game.GetReplyExchangeName())
	status = gr.broker.ConnectQueue(ctx, game.GetReplyExchangeName(), game.GetReplyQueueName(), game.GetReplyKey())
	if status != ports.MESSAGE_SERVICE_OK {
		gr.destructComms(ctx, allocatedQueues, allocatedExchanges, allocatedConnections)
		return GAME_SERVICE_FAILED
	}
	// allocatedConnections = append(allocatedConnections, []string{game.GetReplyExchangeName(), game.GetReplyQueueName(), game})
	return GAME_SERVICE_OK
}
func (gr *GameRegistration) destructComms(ctx context.Context, allocatedQueues []string, allocatedExchanges []string, allocatedConn [][]string) {
	gr.logger.Printf("[destructComms]: Queues to destroy: %s\n", allocatedQueues)
	gr.logger.Printf("[destructComms]: Exchanges to destroy: %s\n", allocatedExchanges)
	gr.logger.Printf("[destructComms]: Connections to destroy: %s\n", allocatedConn)
	for _, connection := range allocatedConn {
		gr.broker.DisconnectQueue(ctx, connection[0], connection[1], connection[2])
	}
	for _, queue := range allocatedQueues {
		gr.broker.DeleteQueue(ctx, queue)
	}
	for _, exchange := range allocatedExchanges {
		gr.broker.DeleteExchange(ctx, exchange)
	}

}

func (gr *GameRegistration) connectStudent(ctx context.Context, game *aggregate.Game, studentID string) GameServiceStatus {
	sQueueName := game.GetStudentControlQueueName(studentID)
	status := gr.broker.CreateQueue(ctx, sQueueName, ports.QUEUE_NO_REPLAY)
	if status != ports.MESSAGE_SERVICE_OK {
		return GAME_SERVICE_FAILED
	}
	status = gr.broker.ConnectQueue(ctx, game.GetControlExchangeName(), sQueueName, game.GetControlKey())
	if status != ports.MESSAGE_SERVICE_OK {
		gr.broker.DeleteQueue(ctx, sQueueName)
		return GAME_SERVICE_FAILED
	}
	return GAME_SERVICE_OK
}

func (gr *GameRegistration) CreateGame(ctx context.Context, requestDTO DTO.CreateGameRequest) (DTO.CreateGameResponse, GameServiceStatus) {

	gameresponse := DTO.CreateGameResponse{}
	gr.logger.Printf("[CreateGame]: Creating Game %s\n", requestDTO.TeacherID)
	game, gamestatus := aggregate.NewGame(requestDTO.TeacherID)
	gr.logger.Printf("[CreateGame]: Creating Game %s\n", game.GetID())
	if gamestatus != aggregate.GAME_OK {
		return DTO.CreateGameResponse{}, GAME_SERVICE_FAILED
	}

	gid, repostatus := gr.gamerepo.CreateGame(ctx, *game)

	gameresponse.GameID = gid
	gameresponse.MovesStream = game.GetMoveQueueName()
	gameresponse.MovesKey = game.GetMovesKey()
	gameresponse.ControlExchange = game.GetControlExchangeName()
	gameresponse.ResponseQueue = game.GetReplyQueueName()
	gameresponse.ControlKey = game.GetControlKey()

	if repostatus == ports.GAMEREPO_ENTRY_EXISTS {
		return gameresponse, GAME_SERVICE_OK
	}
	status := gr.constructComms(ctx, game)
	if status != GAME_SERVICE_OK {
		gr.gamerepo.DeleteGame(ctx, gid)
		return DTO.CreateGameResponse{}, GAME_SERVICE_FAILED
	}
	return gameresponse, GAME_SERVICE_OK
}

func (gr *GameRegistration) ConnectGame(ctx context.Context, requestDTO DTO.ConnectGameRequest) (DTO.ConnectGameResponse, GameServiceStatus) {
	response := DTO.ConnectGameResponse{}
	game, repostatus := gr.gamerepo.GetGame(ctx, requestDTO.GameID)
	if repostatus != ports.GAMEREPO_OK {
		return DTO.ConnectGameResponse{}, GAME_SERVICE_FAILED
	}
	response.GameID = game.GetID()
	response.ControlQueue = game.GetStudentControlQueueName(requestDTO.StudentID)
	response.MovesQueue = game.GetMoveQueueName()
	response.ResponseExchange = game.GetReplyExchangeName()
	response.ResponseKey = game.GetReplyKey()
	gamestatus := game.AddStudent(requestDTO.StudentID)
	if gamestatus == aggregate.STUDENT_EXISTS {
		return response, GAME_SERVICE_OK
	}
	if gamestatus == aggregate.GAME_FAILED {

		return DTO.ConnectGameResponse{}, GAME_SERVICE_FAILED
	}
	status := gr.connectStudent(ctx, &game, requestDTO.StudentID)
	if status == GAME_SERVICE_FAILED {
		return DTO.ConnectGameResponse{}, GAME_SERVICE_FAILED
	}
	gr.gamerepo.SetGame(ctx, game.GetID(), game)
	return response, GAME_SERVICE_OK

}

func (gr *GameRegistration) DeleteGame(ctx context.Context, requestDTO DTO.DeleteGameRequest) GameServiceStatus {
	game, repostatus := gr.gamerepo.GetGame(ctx, requestDTO.GameID)
	if repostatus != ports.GAMEREPO_OK {
		return GAME_SERVICE_FAILED
	}
	allocatedConn := make([][]string, 0)
	allocatedExchanges := make([]string, 0)
	allocatedQueue := make([]string, 0)
	allocatedExchanges = append(allocatedExchanges, game.GetControlExchangeName())
	allocatedExchanges = append(allocatedExchanges, game.GetReplyExchangeName())
	allocatedQueue = append(allocatedQueue, game.GetMoveQueueName())
	allocatedQueue = append(allocatedQueue, game.GetReplyQueueName())
	sIDs := game.GetStudentIDs()
	gr.logger.Printf("Printing Students %s\n", sIDs)
	for _, sID := range sIDs {

		sQueue := game.GetStudentControlQueueName(sID)
		allocatedQueue = append(allocatedQueue, sQueue)
		allocatedConn = append(allocatedConn, []string{game.GetControlExchangeName(), sQueue, game.GetControlKey()})
	}
	allocatedConn = append(allocatedConn, []string{game.GetReplyExchangeName(), game.GetReplyQueueName(), game.GetReplyKey()})
	gr.destructComms(ctx, allocatedQueue, allocatedExchanges, allocatedConn)
	gr.gamerepo.DeleteGame(ctx, requestDTO.GameID)
	return GAME_SERVICE_OK
}
