package pdfgen

import (
	"fmt"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/constants/location"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"github.com/muchlist/risa_restfull/utils/timegen"
	"os"
	"strconv"
	"strings"
)

type PDFReqMonth struct {
	Name        string
	HistoryList dto.HistoryUnwindResponseList
	Start       int64
	End         int64
}

func GeneratePDFVendorMonthly(
	data PDFReqMonth, cekCCTV []dto.VenPhyCheck, cekAltai []dto.AltaiPhyCheck,
) error {
	// slice yang sudah di filter dan dimodifikasi isinya
	var allListComputed []dto.HistoryUnwindResponse
	var completeList []dto.HistoryUnwindResponse
	var completeListNoMaint []dto.HistoryUnwindResponse
	var progressList []dto.HistoryUnwindResponse
	var pendingList []dto.HistoryUnwindResponse

	// idTemp menyimpan id, karena akan banyak id yang sama, maka akan diambil history yang terakhir
	// urutan unwind dengan asumsi unwind sorted by updates.time 1 (pertama kali update tampil pertama)
	var idTemp string
	for _, history := range data.HistoryList {
		// skip jika waktu updatenya melebihi time end laporan
		if history.Updates.Time > data.End {
			continue
		}

		// blok if yang dijalankan jika historynya sama
		if idTemp == history.ID.Hex() {
			if allListComputed == nil {
				continue
			}
			// menambahkan nama pengupdate
			updatedByExisting := allListComputed[len(allListComputed)-1].UpdatedBy
			updatedByCurrent := strings.Split(history.Updates.UpdatedBy, " ")[0]
			if updatedByExisting != updatedByCurrent {
				allListComputed[len(allListComputed)-1].UpdatedBy = updatedByExisting + " > " + updatedByCurrent
			}

			// menambahkan waktu pengerjaan, jika statusComplete sebelumnya pending maka waktu tidak ditambahkan
			difference := history.Updates.Time - allListComputed[len(allListComputed)-1].Updates.Time
			if allListComputed[len(allListComputed)-1].Updates.CompleteStatus == enum.HPending {
				difference = 0
			}

			timeToConsumeExisting := allListComputed[len(allListComputed)-1].UpdatedAt
			allListComputed[len(allListComputed)-1].UpdatedAt = timeToConsumeExisting + difference
			allListComputed[len(allListComputed)-1].Updates = history.Updates
			continue
		}
		// end blok

		idTemp = history.ID.Hex()

		// updatedAt tidak lagi dipakai pada history versi 2,
		// updatedAt akan dialih fungsikan untuk menghitung seberapa lama pekerjaannya diselesaikan
		// rumus createdAt - updatedAt tidak berlaku karena apabila statusCompleted nya pending tidak boleh dihitung
		// terpaksa menggunakan field bertipe int64 lain untuk menampung perhitungan sementara belum memiliki solusi lain
		// updatedAt di nol kan pada data pertama dan akan ditambah jika ada history yang sama
		history.UpdatedAt = 0
		history.UpdatedBy = strings.Split(history.Updates.UpdatedBy, " ")[0]

		allListComputed = append(allListComputed, history)
	}

	for _, historyComputed := range allListComputed {
		if historyComputed.Updates.CompleteStatus == 0 || historyComputed.Updates.CompleteStatus == 4 {
			if strings.ToUpper(historyComputed.Status) != "MAINTENANCE" {
				completeListNoMaint = append(completeListNoMaint, historyComputed)
			}
			completeList = append(completeList, historyComputed)
		}
		if historyComputed.Updates.CompleteStatus == 1 {
			progressList = append(progressList, historyComputed)
		}
		if historyComputed.Updates.CompleteStatus == 2 || historyComputed.Updates.CompleteStatus == 3 {
			pendingList = append(pendingList, historyComputed)
		}
	}

	m := pdf.NewMaroto(consts.Landscape, consts.A4)
	m.SetPageMargins(5, 10, 5)

	startWita, _ := timegen.GetTimeWithYearWITA(data.Start)
	endWita, _ := timegen.GetTimeWithYearWITA(data.End)
	subtitle := fmt.Sprintf("Tanggal %s sd %s", startWita, endWita)
	err := buildHeadingVendorMonth(m, "Rekap Pekerjaan Bulanan di Regional Kalimantan", subtitle)
	if err != nil {
		return err
	}

	if len(completeList) != 0 {
		buildHistoryVendorListMonth(m, completeList, " Completed", getTealColor())
	}

	if len(progressList) != 0 {
		buildHistoryVendorListMonth(m, progressList, " Progress", getOrangeColor())
	}

	if len(pendingList) != 0 {
		buildHistoryVendorListMonth(m, pendingList, " Pending", getPinkColor())
	}

	if len(completeListNoMaint) != 0 {
		m.AddPage()
		err := buildHeadingVendorMonth(m, "Rekap Insiden di Regional Kalimantan", subtitle)
		if err != nil {
			return err
		}
		buildHistoryVendorIncident(m, completeListNoMaint, " Insiden", getTealColor())
	}

	m.AddPage()
	buildCekVendorListMonth(m, cekCCTV, cekAltai, " Cek Bulanan", getPinkColor())

	err = m.OutputFileAndClose(fmt.Sprintf("static/pdf-v-month/%s.pdf", data.Name))
	if err != nil {
		return err
	}
	return nil
}

func buildHeadingVendorMonth(m pdf.Maroto, title string, subtitle string) error {
	var errTemp error
	//m.Row(10, func() {
	//
	//})
	m.Row(20, func() {
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
		m.Col(8, func() {
			textH1(m, title)
			textBodyCenter(m, subtitle, 12)
		})
		m.ColSpace(2)
	})
	return errTemp
}

func buildHistoryVendorListMonth(m pdf.Maroto, dataList []dto.HistoryUnwindResponse, title string, customColor color.Color) {
	tableHeading := []string{"No.", "Nama", "Keterangan", "Solusi", "Pengerjaan", "Update", "Oleh"}

	var contents [][]string
	for i, data := range dataList {
		updateAt, err := timegen.GetTimeWITA(data.Updates.Time)
		if err != nil {
			updateAt = "error"
		}

		contents = append(contents, []string{
			fmt.Sprintf("%03d\n", i+1),
			data.ParentName,
			data.Updates.Problem,
			data.Updates.ProblemResolve,
			sfunc.IntToTime(data.UpdatedAt, ""), // data UpdatedAt sudah diubah pada komputasi sebelumnya menjadi lama pengerjaan
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
			GridSizes: []uint{1, 2, 3, 3, 1, 1, 1},
		},
		ContentProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{1, 2, 3, 3, 1, 1, 1},
		},
		Align:                consts.Left,
		AlternatedBackground: &lightPurpleColor,
		HeaderContentSpace:   1,
		Line:                 true,
	})
}

func buildCekVendorListMonth(m pdf.Maroto, cctvCekList []dto.VenPhyCheck, altaiCekList []dto.AltaiPhyCheck, title string, customColor color.Color) {
	tableHeading := []string{"Nama", "Jenis", "Total", "Dicek", "Belum Cek", "Di MT", "Belum MT", "Dimulai", "Selesai"}

	// 2 , 1,1,  1, 1, 1, 1, 2, 2
	var contents [][]string
	for _, data := range cctvCekList {
		counts := countCekCctv(data)
		updatedAt, _ := timegen.GetTimeWithYearWITA(data.TimeEnded)
		startAt, _ := timegen.GetTimeWithYearWITA(data.TimeStarted)

		contents = append(contents, []string{
			data.Name,
			"CCTV",
			strconv.Itoa(counts.total),
			strconv.Itoa(counts.checked),
			strconv.Itoa(counts.notChecked),
			strconv.Itoa(counts.maintained),
			strconv.Itoa(counts.notMaintained),
			startAt,
			updatedAt},
		)
	}

	for _, data := range altaiCekList {
		counts := countCekAltai(data)
		updatedAt, _ := timegen.GetTimeWithYearWITA(data.TimeEnded)
		startAt, _ := timegen.GetTimeWithYearWITA(data.TimeStarted)

		contents = append(contents, []string{
			data.Name,
			"ALTAI",
			strconv.Itoa(counts.total),
			strconv.Itoa(counts.checked),
			strconv.Itoa(counts.notChecked),
			strconv.Itoa(counts.maintained),
			strconv.Itoa(counts.notMaintained),
			startAt,
			updatedAt},
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
			GridSizes: []uint{2, 1, 1, 1, 1, 1, 1, 2, 2},
		},
		ContentProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{2, 1, 1, 1, 1, 1, 1, 2, 2},
		},
		Align:                consts.Left,
		AlternatedBackground: &lightPurpleColor,
		HeaderContentSpace:   1,
		Line:                 true,
	})
}

type cekData struct {
	total         int
	checked       int
	notChecked    int
	maintained    int
	notMaintained int
}

func countCekCctv(data dto.VenPhyCheck) cekData {
	total := 0
	checked := 0
	notChecked := 0
	maintained := 0
	notMaintained := 0

	for _, cctv := range data.VenPhyCheckItems {
		if !data.QuarterlyMode {
			if cctv.Location == location.Pulpis {
				continue
			}
		}
		total++
		if cctv.IsChecked {
			checked++
		} else {
			notChecked++
		}
		if cctv.IsMaintained {
			maintained++
		} else {
			notMaintained++
		}
	}

	return cekData{
		total:         total,
		checked:       checked,
		notChecked:    notChecked,
		maintained:    maintained,
		notMaintained: notMaintained,
	}
}

func countCekAltai(data dto.AltaiPhyCheck) cekData {
	checked := 0
	notChecked := 0
	maintained := 0
	notMaintained := 0

	for _, altai := range data.AltaiPhyCheckItems {
		if altai.IsChecked {
			checked++
		} else {
			notChecked++
		}
		if altai.IsMaintained {
			maintained++
		} else {
			notMaintained++
		}
	}

	return cekData{
		total:         len(data.AltaiPhyCheckItems),
		checked:       checked,
		notChecked:    notChecked,
		maintained:    maintained,
		notMaintained: notMaintained,
	}
}

func buildHistoryVendorIncident(m pdf.Maroto, dataList []dto.HistoryUnwindResponse, title string, customColor color.Color) {
	//tableHeading := []string{"No.", "Nama", "Keterangan", "Solusi", "Pengerjaan", "Update", "Oleh"}
	lightPurpleColor := getLightPurpleColor()
	m.SetBackgroundColor(customColor)
	m.Row(12, func() {
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

	//GridSizes: []uint{1, 2, 1, 3, 3, 1, 1},

	addListHeader(m)

	for i, data := range dataList {
		if i%2 == 0 {
			m.SetBackgroundColor(lightPurpleColor)
		}
		addList(m, data, i+1)
		m.SetBackgroundColor(color.NewWhite())

	}
}

func addListHeader(m pdf.Maroto) {
	m.Row(9, func() {

		// nomor
		m.Col(1, func() {
			textListB(m, "No.", 0)
		})

		// nama
		m.Col(1, func() {
			textListB(m, "Nama", 0)
		})
		// keterangan
		m.Col(3, func() {
			textListB(m, "Keterangan", 0)
		})

		// solusi
		m.Col(3, func() {
			textListB(m, "Solusi", 0)
		})

		// pengerjaan
		m.Col(1, func() {
			textListB(m, "Pengerjaan", 0)
		})

		// update
		m.Col(1, func() {
			textListB(m, "Update", 0)
		})

		// gambar
		m.Col(2, func() {
			textListB(m, "Gambar", 0)
		})

	})
}

func addList(m pdf.Maroto, data dto.HistoryUnwindResponse, i int) {
	imageFile := ""
	if data.Image != "" {
		split := strings.Split(data.Image, "/")
		imageFile = fmt.Sprintf("static/image/history/thumb_%s", split[len(split)-1])
		if _, err := os.Stat(imageFile); os.IsNotExist(err) {
			imageFile = ""
		}
	}

	updateAt, err := timegen.GetTimeWITA(data.Updates.Time)
	if err != nil {
		updateAt = "error"
	}

	m.Row(20, func() {

		// no
		m.Col(1, func() {
			textList(m, fmt.Sprintf("%03d\n", i), 0)
		})

		// nama
		m.Col(1, func() {
			textList(m, data.ParentName, 0)
		})
		// keterangan
		m.Col(3, func() {
			textList(m, data.Problem, 0)
		})

		// solusi
		m.Col(3, func() {
			textList(m, data.ProblemResolve, 0)
		})

		// pengerjaan
		m.Col(1, func() {
			textList(m, sfunc.IntToTime(data.UpdatedAt, ""), 0)
		})

		// update
		m.Col(1, func() {
			textList(m, updateAt, 0)
		})

		// gambar
		if imageFile != "" {
			m.Col(2, func() {
				_ = m.FileImage(imageFile, props.Rect{
					Left:    0,
					Top:     0,
					Percent: 100,
					Center:  false,
				})
			})
		} else {
			m.Col(2, func() {

			})
		}

	})
}
