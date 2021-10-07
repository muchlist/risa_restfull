package pendingreportdao

import (
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddParticipantInput struct {
	ID           primitive.ObjectID
	Participant  dto.Participant
	FilterBranch string
	UpdatedAt    int64
	UpdatedBy    string
	UpdatedByID  string
}
