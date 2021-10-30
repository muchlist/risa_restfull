package pdfgen

import (
	"github.com/muchlist/risa_restfull/constants/location"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/timegen"
	"strconv"
	"strings"
)

const (
	cctvOfflineKey  = "CCTV Offline"
	cctvBlurKey     = "CCTV Blur"
	altaiOfflineKey = "Altai Offline"
)

type cctvDailyData struct {
	created string
	total   string
	checked string
	blur    string
	offline string
	ok      string
}

type altaiDailyData struct {
	created string
	total   string
	checked string
	offline string
	ok      string
}

type virtualTrouble struct {
	category string
	item     string
}

func convertDailyToDailyViewData(cctv *dto.VendorCheck, altai *dto.AltaiCheck) (cctvDailyData, altaiDailyData, []virtualTrouble) {
	resCctv := cctvDailyData{}
	resAltai := altaiDailyData{}
	deviceTroubleData := make([]virtualTrouble, 0)
	deviceProblem := map[string]string{}

	if cctv != nil {
		checkedTemp := 0
		blurTemp := 0
		offlineTemp := 0
		okTemp := 0
		total := 0

		for _, check := range cctv.VendorCheckItems {
			// jika DisVendor maka kecualikan dari laporan
			if check.DisVendor {
				continue
			}
			if check.IsChecked {
				checkedTemp++
			}
			if check.IsBlur {
				blurTemp++
				deviceProblem[cctvBlurKey] = deviceProblem[cctvBlurKey] + check.Name + ", "
			}
			if check.IsOffline {
				offlineTemp++
				deviceProblem[cctvOfflineKey] = deviceProblem[cctvOfflineKey] + check.Name + ", "
			}
			if !check.IsOffline && !check.IsBlur {
				okTemp++
			}
			total++
		}

		resCctv.created, _ = timegen.GetTimeWithYearWITA(cctv.CreatedAt)
		resCctv.total = strconv.Itoa(total)
		resCctv.ok = strconv.Itoa(okTemp)
		resCctv.checked = strconv.Itoa(checkedTemp)
		resCctv.blur = strconv.Itoa(blurTemp)
		resCctv.offline = strconv.Itoa(offlineTemp)
	}

	if altai != nil {
		checkedTemp := 0
		offlineTemp := 0
		okTemp := 0
		total := 0

		for _, check := range altai.AltaiCheckItems {
			// jika DisVendor maka kecualikan dari laporan
			if check.DisVendor {
				continue
			}
			if check.IsChecked {
				checkedTemp++
			}
			if check.IsOffline {
				offlineTemp++
				deviceProblem[altaiOfflineKey] = deviceProblem[altaiOfflineKey] + check.Name + ", "
			}
			if !check.IsOffline {
				okTemp++
			}
			total++
		}

		resAltai.created, _ = timegen.GetTimeWithYearWITA(altai.CreatedAt)
		resAltai.total = strconv.Itoa(total)
		resAltai.ok = strconv.Itoa(okTemp)
		resAltai.checked = strconv.Itoa(checkedTemp)
		resAltai.offline = strconv.Itoa(offlineTemp)
	}

	for key, element := range deviceProblem {
		deviceTroubleData = append(deviceTroubleData, virtualTrouble{
			category: key,
			item:     strings.TrimSuffix(element, ", "),
		})
	}

	return resCctv, resAltai, deviceTroubleData
}

// =======================================================================

type summaryMonthlyData struct {
	created    string
	total      string
	checked    string
	notChecked string
}

func convertMonthlyViewData(cctv *dto.VenPhyCheck, altai *dto.AltaiPhyCheck) (cctvRes summaryMonthlyData, altaiRes summaryMonthlyData) {
	if cctv != nil {
		checkedTemp := 0
		totalTemp := 0

		for _, check := range cctv.VenPhyCheckItems {
			if check.DisVendor {
				continue
			}
			if check.Location != location.Pulpis {
				totalTemp++
			}
			if check.IsChecked {
				checkedTemp++
			}
		}

		cctvRes.created, _ = timegen.GetTimeWithYearWITA(cctv.CreatedAt)
		cctvRes.total = strconv.Itoa(totalTemp)
		cctvRes.checked = strconv.Itoa(checkedTemp)
		cctvRes.notChecked = strconv.Itoa(totalTemp - checkedTemp)
	}

	if altai != nil {
		checkedTemp := 0
		totalTemp := 0

		for _, check := range altai.AltaiPhyCheckItems {
			if check.DisVendor {
				continue
			}
			if check.IsChecked {
				checkedTemp++
			}
			totalTemp++
		}

		altaiRes.created, _ = timegen.GetTimeWithYearWITA(altai.CreatedAt)
		altaiRes.total = strconv.Itoa(totalTemp)
		altaiRes.checked = strconv.Itoa(checkedTemp)
		altaiRes.notChecked = strconv.Itoa(totalTemp - checkedTemp)
	}

	return
}

//========================================================================================================

type summaryQuarterlyData struct {
	created       string
	total         string
	maintained    string
	notMaintained string
}

func convertQuarterlyViewDataCctv(cctv *dto.VenPhyCheck) (cctvReg summaryQuarterlyData, cctvPulpis summaryQuarterlyData) {
	if cctv != nil {
		regTotal := 0
		pulpisTotal := 0
		regMaintTemp := 0
		pulpisMaintTemp := 0

		for _, check := range cctv.VenPhyCheckItems {
			if check.DisVendor {
				continue
			}
			// cek puplis
			if check.Location == location.Pulpis {
				pulpisTotal++
				if check.IsMaintained {
					pulpisMaintTemp++
				}
			} else { //  cek selain puplis
				regTotal++
				if check.IsMaintained {
					regMaintTemp++
				}
			}
		}

		cctvReg.created, _ = timegen.GetTimeWithYearWITA(cctv.CreatedAt)
		cctvReg.total = strconv.Itoa(regTotal)
		cctvReg.maintained = strconv.Itoa(regMaintTemp)
		cctvReg.notMaintained = strconv.Itoa(regTotal - regMaintTemp)

		cctvPulpis.created, _ = timegen.GetTimeWithYearWITA(cctv.CreatedAt)
		cctvPulpis.total = strconv.Itoa(pulpisTotal)
		cctvPulpis.maintained = strconv.Itoa(pulpisMaintTemp)
		cctvPulpis.notMaintained = strconv.Itoa(pulpisTotal - pulpisMaintTemp)
	}
	return
}

func convertQuarterlyViewDataAltai(altai *dto.AltaiPhyCheck) summaryQuarterlyData {
	altaiRes := summaryQuarterlyData{}

	if altai != nil {
		maintainTemp := 0
		total := 0

		for _, check := range altai.AltaiPhyCheckItems {
			if check.DisVendor {
				continue
			}
			if check.IsMaintained {
				maintainTemp++
			}
			total++
		}

		altaiRes.created, _ = timegen.GetTimeWithYearWITA(altai.CreatedAt)
		altaiRes.total = strconv.Itoa(total)
		altaiRes.maintained = strconv.Itoa(maintainTemp)
		altaiRes.notMaintained = strconv.Itoa(total - maintainTemp)
	}

	return altaiRes
}
