package pendingreportdao

import (
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ParticipantParams struct {
	ID           primitive.ObjectID
	Participant  dto.Participant
	FilterBranch string
	UpdatedAt    int64
	UpdatedBy    string
	UpdatedByID  string
}

type EditParticipantParams struct {
	ID           primitive.ObjectID
	FilterBranch string
	Participant  []dto.Participant
	Approver     []dto.Participant
	UpdatedAt    int64
	UpdatedBy    string
	UpdatedByID  string
}
