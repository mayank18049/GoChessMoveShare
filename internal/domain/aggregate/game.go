package aggregate

import (
	"github.com/google/uuid"
	"github.com/mayank18049/GoChessMoveShare/internal/domain/model"
)

type Game struct {
	id       string
	teacher  *model.User
	students model.Users
}

func NewGame(teacherID string) (*Game, GameStatus) {
	teacher, err := model.NewUser(teacherID)
	switch err {
	case model.INVALID_USERID:
		return nil, INVALID_TEACHERID
	case model.USER_OK:
		game := Game{
			id:       uuid.New().String(),
			teacher:  teacher,
			students: make([]*model.User, 0),
		}
		return &game, GAME_OK
	default:
		return nil, GAME_FAILED
	}

}

func (g *Game) AddStudent(studentID string) GameStatus {
	student, err := model.NewUser(studentID)
	if err != model.USER_OK {
		return GAME_FAILED
	}
	userstatus := g.students.AddUser(student)
	switch userstatus {
	case model.USER_OK:
		return GAME_OK
	default:
		return STUDENT_EXISTS
	}
}

func (g Game) GetMoveQueueName() string {
	return g.id + "-" + g.teacher.GetID() + "-moves-stream"
}

func (g *Game) GetControlExchangeName() string {
	return g.id + "-" + g.teacher.GetID() + "-control-exchange"
}
func (g *Game) GetReplyExchangeName() string {
	return g.id + "-" + g.teacher.GetID() + "-reply-exchange"
}

func (g *Game) GetReplyQueueName() string {
	return g.id + "-" + g.teacher.GetID() + "-reply-queue"
}

func (g *Game) GetControlKey() string {
	return "control-data"
}
func (g *Game) GetReplyKey() string {
	return "reply"
}
func (g *Game) GetMovesKey() string {
	return "moves"
}

func (g *Game) GetStudentControlQueueName(studentID string) string {
	if g.students.ContainsUserID(studentID) {
		return g.id + "-" + studentID + "-control-queue"
	}
	return ""
}

func (g Game) GetID() string {
	return g.id
}
func (g Game) GetTeacherID() string {
	return g.teacher.GetID()
}

func (g Game) GetStudentIDs() []string {
	sIDs := make([]string, 0)
	for _, student := range g.students {
		sIDs = append(sIDs, student.GetID())
	}
	return sIDs
}
