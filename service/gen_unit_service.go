package service

import (
	"fmt"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/clients/fcm"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/userdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"net"
	"strings"
)

func NewGenUnitService(
	dao genunitdao.GenUnitDaoAssumer,
	daoU userdao.UserDaoAssumer,
	fcmClient fcm.ClientAssumer) GenUnitServiceAssumer {
	return &genUnitService{
		dao:       dao,
		daoU:      daoU,
		fcmClient: fcmClient,
	}
}

type genUnitService struct {
	dao       genunitdao.GenUnitDaoAssumer
	daoU      userdao.UserDaoAssumer
	fcmClient fcm.ClientAssumer
}

type GenUnitServiceAssumer interface {
	FindUnit(filter dto.GenUnitFilter) (dto.GenUnitResponseList, rest_err.APIError)
	GetIPList(branchIfSpecific string, category string) ([]string, rest_err.APIError)
	AppendPingState(input dto.GenUnitPingStateRequest) (int64, rest_err.APIError)
	CheckHardwareDownAndSendNotif(branchIfSpecific string, category string) rest_err.APIError
}

func (g *genUnitService) FindUnit(filter dto.GenUnitFilter) (dto.GenUnitResponseList, rest_err.APIError) {
	// cek apakah ip address valid, jika valid maka set filter.FilterName ke kosong supaya pencarian berdasarkan FilterIP
	if filter.IP != "" {
		if net.ParseIP(filter.IP) == nil {
			return nil, rest_err.NewBadRequestError("IP Address tidak valid")
		}
		filter.Name = ""
	}

	// DB
	unitList, err := g.dao.FindUnit(filter)
	if err != nil {
		return nil, err
	}
	return unitList, nil
}

func (g *genUnitService) GetIPList(branchIfSpecific string, category string) ([]string, rest_err.APIError) {
	// DB
	ipAddressList, err := g.dao.GetIPList(branchIfSpecific, category)
	if err != nil {
		return nil, err
	}

	uniqueIPList := sfunc.Unique(ipAddressList)
	return uniqueIPList, nil
}

// CheckHardwareDownAndSendNotif mengecek perangkat yang mati dan belum pernah dilakukan pengeckean
// lalu mengirimkan notif ke semua user terkait
// sudah dilengkapi dengan logger
func (g *genUnitService) CheckHardwareDownAndSendNotif(branchIfSpecific string, category string) rest_err.APIError {
	// DB
	unitList, err := g.dao.FindUnit(dto.GenUnitFilter{
		Branch:   branchIfSpecific,
		Category: category,
		Pings:    true,
		LastPing: "DOWN",
	})
	if err != nil {
		logger.Error("mendapatkan unit gagal saat menambahkan fcm (CheckHardwareDownAndSendNotif)", err)
		return err
	}

	// filtering
	filterUncheckedUnitList(&unitList)

	users, err := g.daoU.FindUser(branchIfSpecific)
	if err != nil {
		logger.Error("mendapatkan user gagal saat menambahkan fcm (CheckHardwareDownAndSendNotif)", err)
		return err
	}
	var tokens []string
	for _, u := range users {
		tokens = append(tokens, u.FcmToken)
	}

	var allUnitBuilder strings.Builder
	for _, unit := range unitList {
		allUnitBuilder.WriteString(unit.Name + ", ")
	}
	allUnitString := strings.TrimRight(allUnitBuilder.String(), ", ")

	// firebase
	g.fcmClient.SendMessage(fcm.Payload{
		Title:          fmt.Sprintf("%d unit perlu dicek", len(unitList)),
		Message:        fmt.Sprint(allUnitString),
		ReceiverTokens: tokens,
	})

	return err
}

// hanya sisakan unit yang tidak memiliki case dan sedang down total.
func filterUncheckedUnitList(data *dto.GenUnitResponseList) {
	temp := dto.GenUnitResponseList{}
	for _, gen := range *data {
		if len(gen.PingsState) == 0 {
			continue
		}
		var totalPing int
		for _, ping := range gen.PingsState {
			totalPing = +ping.Code
		}
		// cek casenya harus 0
		isCaseZero := gen.CasesSize == 0
		isTotalPingZero := totalPing == 0

		if isCaseZero && isTotalPingZero {
			temp = append(temp, gen)
		}
	}
	*data = temp
}

func (g *genUnitService) AppendPingState(input dto.GenUnitPingStateRequest) (int64, rest_err.APIError) {
	// DB
	unitUpdatedCount, err := g.dao.AppendPingState(input)
	if err != nil {
		return 0, err
	}

	return unitUpdatedCount, nil
}
