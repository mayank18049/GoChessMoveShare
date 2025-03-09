package aggregate

type GameStatus int
type Games []Game

const (
	GAME_OK GameStatus = -iota
	STUDENT_NOT_FOUND
	TEAHCER_EXISTS
	STUDENT_EXISTS
	INVALID_STUDENTID
	INVALID_TEACHERID
	INVALID_GAMEID
	GAME_FAILED
)
