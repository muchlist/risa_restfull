package pdfgen

import (
	"fmt"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/timegen"
	"strconv"
	"strings"
)

type PDFReq struct {
	Name        string
	HistoryList dto.HistoryUnwindResponseList
	CheckList   []dto.Check
	Start       int64
	End         int64
}

func GeneratePDF(
	pdfStruct PDFReq,
) error {

	// slice yang sudah di filter dan dimodifikasi isinya
	var allListComputed []dto.HistoryUnwindResponse

	var completeList []dto.HistoryUnwindResponse
	var progressList []dto.HistoryUnwindResponse
	var pendingList []dto.HistoryUnwindResponse

	// idTemp menyimpan id, karena akan banyak id yang sama, maka akan diambil history yang terakhir
	// urutan unwind dengan asumsi unwind sorted by updates.time 1 (pertama kali update tampil pertama)
	var idTemp string
	for _, history := range pdfStruct.HistoryList {
		// skip jika waktu updatenya melebihi time end laporan
		if history.Updates.Time > pdfStruct.End {
			continue
		}

		if strings.ToUpper(history.Status) == "MAINTENANCE" {
			continue
		}

		// blok if yang dijalankan jika historynya sama
		if idTemp == history.ID.Hex() {
			if allListComputed != nil {
				updatedByExisting := allListComputed[len(allListComputed)-1].UpdatedBy
				updatedByCurrent := strings.Split(history.Updates.UpdatedBy, " ")[0]
				if updatedByExisting != updatedByCurrent {
					allListComputed[len(allListComputed)-1].UpdatedBy = updatedByExisting + " > " + updatedByCurrent
				}
				allListComputed[len(allListComputed)-1].Updates = history.Updates
				continue
			}
		}
		// end blok

		idTemp = history.ID.Hex()

		history.UpdatedBy = strings.Split(history.Updates.UpdatedBy, " ")[0]

		allListComputed = append(allListComputed, history)
	}

	for _, historyComputed := range allListComputed {
		if historyComputed.Updates.CompleteStatus == enum.HInfo ||
			historyComputed.Updates.CompleteStatus == enum.HComplete ||
			historyComputed.Updates.CompleteStatus == enum.HRequestComplete ||
			historyComputed.Updates.CompleteStatus == enum.HCompleteWithBA {
			completeList = append(completeList, historyComputed)
		}
		if historyComputed.Updates.CompleteStatus == enum.HProgress {
			progressList = append(progressList, historyComputed)
		}
		if historyComputed.Updates.CompleteStatus == enum.HRequestPending ||
			historyComputed.Updates.CompleteStatus == enum.HPending {
			pendingList = append(pendingList, historyComputed)
		}
	}

	m := pdf.NewMaroto(consts.Landscape, consts.A4)
	m.SetPageMargins(5, 10, 5)

	startWita, _ := timegen.GetTimeWithYearWITA(pdfStruct.Start)
	endWita, _ := timegen.GetTimeWithYearWITA(pdfStruct.End)
	subtitle := fmt.Sprintf("Tanggal %s sd %s", startWita, endWita)
	err := buildHeading(m, subtitle)
	if err != nil {
		return err
	}

	if len(completeList) != 0 {
		buildHistoryList(m, completeList, " Completed", getTealColor())
	}

	if len(progressList) != 0 {
		buildHistoryList(m, progressList, " Progress", getOrangeColor())
	}

	if len(pendingList) != 0 {
		buildHistoryList(m, pendingList, " Pending", getPinkColor())
	}

	if len(pdfStruct.CheckList) != 0 {
		buildCheckList(m, pdfStruct.CheckList)
	}

	err = m.OutputFileAndClose(fmt.Sprintf("static/pdf/%s.pdf", pdfStruct.Name))
	if err != nil {
		return err
	}
	return nil
}

func buildHeading(m pdf.Maroto, subtitle string) error {
	var errTemp error
	m.Row(10, func() {

	})
	m.Row(20, func() {
		m.ColSpace(2)
		m.Col(8, func() {
			textH1(m, "Rekap Laporan IT Regional Kalimantan")
			textBodyCenter(m, subtitle, 12)
		})
		m.Col(2, func() {
			err := m.FileImage("static/image/pelindo3.png", props.Rect{
				Percent: 100,
				Center:  false,
				Top:     3,
			})
			if err != nil {
				errTemp = err
			}
		})
	})
	return errTemp
}

func buildHistoryList(m pdf.Maroto, dataList []dto.HistoryUnwindResponse, title string, customColor color.Color) {
	tableHeading := []string{"Nama", "Kategori", "Keterangan", "Solusi", "Status", "Update", "Oleh"}

	var contents [][]string
	for _, data := range dataList {
		updateAt, err := timegen.GetTimeWITA(data.Updates.Time)
		if err != nil {
			updateAt = "error"
		}

		contents = append(contents, []string{
			data.ParentName,
			data.Category,
			data.Updates.Problem,
			data.Updates.ProblemResolve,
			enum.GetProgressString(data.Updates.CompleteStatus),
			updateAt,
			strings.ToLower(data.UpdatedBy)},
		)
	}

	lightPurpleColor := getLightPurpleColor()

	m.SetBackgroundColor(customColor)
	m.Row(9, func() {
		m.Col(12, func() {
			m.Text(title, props.Text{
				Top:             2,
				Family:          consts.Courier,
				Style:           consts.Bold,
				Size:            12,
				Align:           consts.Left,
				VerticalPadding: 0,
				Color:           color.NewWhite(),
			})
		})
	})
	m.SetBackgroundColor(color.NewWhite())
	m.TableList(tableHeading, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{2, 1, 3, 3, 1, 1, 1},
		},
		ContentProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{2, 1, 3, 3, 1, 1, 1},
		},
		Align:                consts.Left,
		AlternatedBackground: &lightPurpleColor,
		HeaderContentSpace:   1,
		Line:                 true,
	})
}

func buildCheckList(m pdf.Maroto, checkList []dto.Check) {
	tableHeading := []string{"Judul", "Shift", "Lokasi", "Keterangan", "Problem", "Cek", "Oleh"}

	var contents [][]string

	for _, check := range checkList {
		for _, data := range check.CheckItems {
			checkedAt, err := timegen.GetHourWITA(data.CheckedAt)
			if err != nil {
				checkedAt = "error"
			}
			if data.CheckedAt == 0 {
				checkedAt = "tidak dicek"
				data.CheckedNote = ""
			}

			haveProblem := ""
			if data.HaveProblem {
				haveProblem = "ada"
			}

			contents = append(contents, []string{
				data.Name,
				strconv.Itoa(check.Shift),
				data.Location,
				data.CheckedNote,
				haveProblem,
				checkedAt,
				strings.Split(check.CreatedBy, " ")[0]})
		}
	}

	lightPurpleColor := getLightPurpleColor()

	m.SetBackgroundColor(getTealColor())
	m.Row(9, func() {
		m.Col(12, func() {
			m.Text(" CheckList", props.Text{
				Top:             2,
				Family:          consts.Courier,
				Style:           consts.Bold,
				Size:            12,
				Align:           consts.Left,
				VerticalPadding: 0,
				Color:           color.NewWhite(),
			})
		})
	})
	m.SetBackgroundColor(color.NewWhite())
	m.TableList(tableHeading, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{3, 1, 1, 4, 1, 1, 1},
		},
		ContentProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{3, 1, 1, 4, 1, 1, 1},
		},
		Align:                consts.Left,
		AlternatedBackground: &lightPurpleColor,
		HeaderContentSpace:   1,
		Line:                 true,
	})
}
