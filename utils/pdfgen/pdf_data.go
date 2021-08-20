package pdfgen

import (
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

		for _, check := range cctv.VendorCheckItems {
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
		}

		resCctv.created, _ = timegen.GetTimeWithYearWITA(cctv.CreatedAt)
		resCctv.total = strconv.Itoa(len(cctv.VendorCheckItems))
		resCctv.ok = strconv.Itoa(okTemp)
		resCctv.checked = strconv.Itoa(checkedTemp)
		resCctv.blur = strconv.Itoa(blurTemp)
		resCctv.offline = strconv.Itoa(offlineTemp)
	}

	if altai != nil {
		checkedTemp := 0
		offlineTemp := 0
		okTemp := 0

		for _, check := range altai.AltaiCheckItems {
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
		}

		resAltai.created, _ = timegen.GetTimeWithYearWITA(altai.CreatedAt)
		resAltai.total = strconv.Itoa(len(altai.AltaiCheckItems))
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

		for _, check := range cctv.VenPhyCheckItems {
			if check.IsChecked {
				checkedTemp++
			}
		}

		cctvRes.created, _ = timegen.GetTimeWithYearWITA(cctv.CreatedAt)
		cctvRes.total = strconv.Itoa(len(cctv.VenPhyCheckItems))
		cctvRes.checked = strconv.Itoa(checkedTemp)
		cctvRes.notChecked = strconv.Itoa(len(cctv.VenPhyCheckItems) - checkedTemp)
	}

	if altai != nil {
		checkedTemp := 0

		for _, check := range altai.AltaiPhyCheckItems {
			if check.IsChecked {
				checkedTemp++
			}
		}

		altaiRes.created, _ = timegen.GetTimeWithYearWITA(altai.CreatedAt)
		altaiRes.total = strconv.Itoa(len(altai.AltaiPhyCheckItems))
		altaiRes.checked = strconv.Itoa(checkedTemp)
		altaiRes.notChecked = strconv.Itoa(len(altai.AltaiPhyCheckItems) - checkedTemp)
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

func convertQuarterlyViewData(cctv *dto.VenPhyCheck, altai *dto.AltaiPhyCheck) (cctvRes summaryQuarterlyData, altaiRes summaryQuarterlyData) {
	if cctv != nil {
		checkedTemp := 0

		for _, check := range cctv.VenPhyCheckItems {
			if check.IsChecked {
				checkedTemp++
			}
		}

		cctvRes.created, _ = timegen.GetTimeWithYearWITA(cctv.CreatedAt)
		cctvRes.total = strconv.Itoa(len(cctv.VenPhyCheckItems))
		cctvRes.maintained = strconv.Itoa(checkedTemp)
		cctvRes.notMaintained = strconv.Itoa(len(cctv.VenPhyCheckItems) - checkedTemp)
	}

	if altai != nil {
		checkedTemp := 0

		for _, check := range altai.AltaiPhyCheckItems {
			if check.IsChecked {
				checkedTemp++
			}
		}

		altaiRes.created, _ = timegen.GetTimeWithYearWITA(altai.CreatedAt)
		altaiRes.total = strconv.Itoa(len(altai.AltaiPhyCheckItems))
		altaiRes.maintained = strconv.Itoa(checkedTemp)
		altaiRes.notMaintained = strconv.Itoa(len(altai.AltaiPhyCheckItems) - checkedTemp)
	}

	return
}
