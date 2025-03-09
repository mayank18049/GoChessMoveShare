package DTO

type CreateGameRequest struct {
	TeacherID string `json:"teacherID"`
}
type CreateGameResponse struct {
	GameID          string `json:"gameID"`
	ControlExchange string `json:"control_exchange"`
	ControlKey      string `json:"control_key"`
	MovesStream     string `json:"move_queue"`
	MovesKey        string `json:"move_key"`
	ResponseQueue   string `json:"response_queue"`
}

type ConnectGameRequest struct {
	GameID    string `json:"gameID"`
	StudentID string `json:"studentID"`
}

type ConnectGameResponse struct {
	GameID           string `json:"gameID"`
	ControlQueue     string `json:"control_queue"`
	ResponseExchange string `json:"response_exchange"`
	ResponseKey      string `json:"response_key"`
	MovesQueue       string `json:"move_queue"`
}

type DeleteGameRequest struct {
	GameID string `json:"gameID"`
}

type DisconnectGameRequest struct {
	GameID    string `json:"gameID"`
	StudentID string `json:"studentID"`
}
