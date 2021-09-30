package pdfgen

import (
	"fmt"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"github.com/muchlist/risa_restfull/utils/timegen"
	"os"
	"strings"
)

type PDFReqMonth struct {
	Name        string
	HistoryList dto.HistoryUnwindResponseList
	Start       int64
	End         int64
}

func GeneratePDFVendorMonthly(
	data PDFReqMonth, dataMaint dto.ReportResponse, dataCheckConf dto.ConfigCheck,
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
		if historyComputed.Updates.CompleteStatus == enum.HInfo ||
			historyComputed.Updates.CompleteStatus == enum.HComplete ||
			historyComputed.Updates.CompleteStatus == enum.HRequestComplete ||
			historyComputed.Updates.CompleteStatus == enum.HCompleteWithBA {
			if strings.ToUpper(historyComputed.Status) != "MAINTENANCE" {
				completeListNoMaint = append(completeListNoMaint, historyComputed)
			}
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
	m.SetPageMargins(10, 10, 10)

	startWita, _ := timegen.GetTimeWithYearWITA(data.Start)
	endWita, _ := timegen.GetTimeWithYearWITA(data.End)
	subtitle := fmt.Sprintf("Tanggal %s sd %s", startWita, endWita)
	err := buildHeadingVendorMonth(m, "Rekap Pekerjaan Bulanan di Regional Kalimantan", subtitle)
	if err != nil {
		return err
	}

	// SPACE
	m.Row(5, func() {
		m.Col(0, func() {
		})
	})

	// MONTHLY
	//----------convert data
	cctvMonthlyViewData, altaiMonthlyViewData := convertMonthlyViewData(dataMaint.CctvMonthly, dataMaint.AltaiMonthly)
	buildTitleHeadingView(m, " Cek Fisik Bulanan", getTealColor())
	buildCCTVMonthlyViewLand(m, cctvMonthlyViewData, altaiMonthlyViewData)

	// SPACE
	//m.Row(5, func() {
	//	m.Col(0, func() {
	//	})
	//})

	// QUARTERLY
	//----------convert data
	regCctvQuarterlyViewData, pulpisCctvQuarterlyViewData := convertQuarterlyViewDataCctv(dataMaint.CctvQuarterly)
	altaiQuarterlyViewData := convertQuarterlyViewDataAltai(dataMaint.AltaiQuarterly)
	buildTitleHeadingView(m, " Cek Fisik Triwulan", getOrangeColor())
	buildCCTVQuarterlyViewLand(m, regCctvQuarterlyViewData, altaiQuarterlyViewData)

	// SPACE
	//m.Row(5, func() {
	//	m.Col(0, func() {
	//	})
	//})

	buildTitleHeadingView(m, " Cek Fisik Triwulan Pulpis", getOrangeColor())
	buildCCTVQuarterlyViewNoAltaiLand(m, pulpisCctvQuarterlyViewData)

	// backup check
	// SPACE
	//m.Row(5, func() {
	//	m.Col(0, func() {
	//	})
	//})
	buildTitleHeadingView(m, " Pencadangan konfigurasi", getDarkGreyColorLight())
	totalConfig := len(dataCheckConf.ConfigCheckItems)
	configUpdated := 0
	for _, d := range dataCheckConf.ConfigCheckItems {
		if d.IsUpdated {
			configUpdated++
		}
	}
	m.Row(5, func() {
		textBody(m, fmt.Sprintf("- Pencadangan konfigurasi : %d dari %d perangkat jaringan dicadangkan (data terlampir)", configUpdated, totalConfig), 5)
	})

	// NEW PAGE ================================================================= \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
	m.AddPage()
	m.SetPageMargins(5, 10, 5)

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

	// NEW PAGE ================================================================= \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
	m.AddPage()
	m.SetPageMargins(10, 10, 10)

	buildConfigListMonth(m, dataCheckConf.ConfigCheckItems, " Lampiran pencadangan konfigurasi perangkat", getDarkGreyColorLight())

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

		data.UpdatedAt = 0 // permintaan bos untuk di 0 kan
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
		m.Col(2, func() {
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
		m.Col(1, func() {
			textListB(m, "Gambar", 0)
		})

	})
}

func addList(m pdf.Maroto, data dto.HistoryUnwindResponse, i int) {
	imageFile := ""
	if data.Image != "" {
		split := strings.Split(data.Image, "/")
		imageFile = fmt.Sprintf("static/image/history/thumb_%s", split[len(split)-1])
		// cek apakah gambar exist
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
		m.Col(2, func() {
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

		data.UpdatedAt = 0 // permintaan bos untuk di 0 kan
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
			m.Col(1, func() {
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

func buildCCTVMonthlyViewLand(m pdf.Maroto, cctv summaryMonthlyData, altai summaryMonthlyData) {

	// CCTV VIRTUAL ----- ALTAI VIRTUAL
	m.Row(15, func() {
		m.Col(5, func() {
			textH3(m, "CCTV", 3)
			textBodyItalic(m, fmt.Sprintf(" dimulai dari %s", cctv.created), 8)
		})
		m.Col(5, func() {
			textH3(m, "ALTAI", 3)
			textBodyItalic(m, fmt.Sprintf(" dimulai dari %s", altai.created), 8)
		})
	})

	// DATA 3-2-1     ------     DATA 3-2-1
	m.Row(20, func() {
		m.Col(3, func() {
			textBody(m, "- Total CCTV", 0)
			textBody(m, "- Sudah di cek", 5)
			textBody(m, "- Belum di cek", 10)
		})

		m.Col(2, func() {
			textBody(m, cctv.total, 0)
			textBody(m, cctv.checked, 5)
			textBody(m, cctv.notChecked, 10)
		})

		m.Col(3, func() {
			textBody(m, "- Total Altai", 0)
			textBody(m, "- Sudah di cek", 5)
			textBody(m, "- Belum di cek", 10)
		})

		m.Col(2, func() {
			textBody(m, altai.total, 0)
			textBody(m, altai.checked, 5)
			textBody(m, altai.notChecked, 10)
		})

	})

}

func buildCCTVQuarterlyViewLand(m pdf.Maroto, cctv summaryQuarterlyData, altai summaryQuarterlyData) {

	// CCTV  ----- ALTAI
	m.Row(15, func() {
		m.Col(5, func() {
			textH3(m, "CCTV", 3)
			textBodyItalic(m, fmt.Sprintf(" dimulai dari %s", cctv.created), 8)
		})
		m.Col(5, func() {
			textH3(m, "ALTAI", 3)
			textBodyItalic(m, fmt.Sprintf(" dimulai dari %s", altai.created), 8)
		})
	})

	// DATA 3-2-1     ------     DATA 3-2-1
	m.Row(20, func() {

		m.Col(3, func() {
			textBody(m, "- Total CCTV", 0)
			textBody(m, "- Sudah di maintenance", 5)
			textBody(m, "- Belum di maintenance", 10)
		})

		m.Col(2, func() {
			textBody(m, cctv.total, 0)
			textBody(m, cctv.maintained, 5)
			textBody(m, cctv.notMaintained, 10)
		})

		m.Col(3, func() {
			textBody(m, "- Total Altai", 0)
			textBody(m, "- Sudah di maintenance", 5)
			textBody(m, "- Belum di maintenance", 10)
		})

		m.Col(2, func() {
			textBody(m, altai.total, 0)
			textBody(m, altai.maintained, 5)
			textBody(m, altai.notMaintained, 10)
		})

	})

}

func buildCCTVQuarterlyViewNoAltaiLand(m pdf.Maroto, cctv summaryQuarterlyData) {

	// CCTV VIRTUAL
	m.Row(15, func() {
		m.Col(5, func() {
			textH3(m, "CCTV", 3)
			textBodyItalic(m, fmt.Sprintf(" dimulai dari %s", cctv.created), 8)
		})
		m.ColSpace(7)
	})

	// DATA 3-2-1
	m.Row(20, func() {
		m.Col(3, func() {
			textBody(m, "- Total CCTV", 0)
			textBody(m, "- Sudah di maintenance", 5)
			textBody(m, "- Belum di maintenance", 10)
		})

		m.Col(2, func() {
			textBody(m, cctv.total, 0)
			textBody(m, cctv.maintained, 5)
			textBody(m, cctv.notMaintained, 10)
		})

		m.ColSpace(7)

	})

}

func buildConfigListMonth(m pdf.Maroto, dataList []dto.ConfigCheckItemEmbed, title string, customColor color.Color) {
	tableHeading := []string{"No.", "Perangkat", "Lokasi", "Status", "Waktu verifikasi", "Diverifikasi oleh"}

	var contents [][]string
	for i, data := range dataList {
		updateAt, err := timegen.GetTimeWITA(data.CheckedAt)
		if err != nil {
			updateAt = "error"
		}

		var status string
		if data.IsUpdated {
			status = "Dicadangkan"
		}

		contents = append(contents, []string{
			fmt.Sprintf("%03d\n", i+1),
			data.Name,
			data.Location,
			status,
			updateAt,
			strings.ToLower(data.CheckedBy)},
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
			GridSizes: []uint{2, 3, 2, 1, 2, 2},
		},
		ContentProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{2, 3, 2, 1, 2, 2},
		},
		Align:                consts.Left,
		AlternatedBackground: &lightPurpleColor,
		HeaderContentSpace:   1,
		Line:                 true,
	})
}
