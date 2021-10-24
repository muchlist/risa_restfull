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
	"strings"
)

func GeneratePDFVendorDaily(name string, data dto.ReportResponse, pdfVendorStruct PDFVendorReq) error {

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 10, 10)

	// HEADING =============================================================================================
	startWita, _ := timegen.GetTimeWithYearWITA(pdfVendorStruct.Start)
	endWita, _ := timegen.GetTimeWithYearWITA(pdfVendorStruct.End)
	subtitle := fmt.Sprintf("Tanggal %s sd %s", startWita, endWita)
	err := buildHeadingVendorDaily(m, subtitle)
	if err != nil {
		return err
	}

	// DAILY CCTV & ALTAI ==================================================================================
	// --- heading
	buildTitleHeadingView(m, " Cek Virtual Harian", getTealColor())
	// ---- body
	// ----------convert data
	cctvDailyViewData, altaiDailyViewData, deviceProblemMap := convertDailyToDailyViewData(data.CctvDaily, data.AltaiDaily)
	buildCCTVDailyView(m, cctvDailyViewData, altaiDailyViewData)
	// ---- rekap perangkat yang perlu ditangani
	if len(deviceProblemMap) != 0 {
		buildVirtualTrouble(m, deviceProblemMap)
	}

	// SPACE
	m.Row(5, func() {
		m.Col(0, func() {
		})
	})

	// MONTHLY
	//----------convert data
	cctvMonthlyViewData, altaiMonthlyViewData := convertMonthlyViewData(data.CctvMonthly, data.AltaiMonthly)
	buildTitleHeadingView(m, " Cek Fisik Bulanan", getOrangeColor())
	buildCCTVMonthlyView(m, cctvMonthlyViewData, altaiMonthlyViewData)

	// SPACE
	m.Row(5, func() {
		m.Col(0, func() {
		})
	})

	// QUARTERLY
	//----------convert data
	regCctvQuarterlyViewData, pulpisCctvQuarterlyViewData := convertQuarterlyViewDataCctv(data.CctvQuarterly)
	altaiQuarterlyViewData := convertQuarterlyViewDataAltai(data.AltaiQuarterly)
	buildTitleHeadingView(m, " Cek Fisik Triwulan", getPinkColor())
	buildCCTVQuarterlyView(m, regCctvQuarterlyViewData, altaiQuarterlyViewData)

	// SPACE
	m.Row(5, func() {
		m.Col(0, func() {
		})
	})

	buildTitleHeadingView(m, " Cek Fisik Triwulan Pulpis", getPinkColor())
	buildCCTVQuarterlyViewNoAltai(m, pulpisCctvQuarterlyViewData)

	// NEW PAGE ================================================================= \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
	m.AddPage()

	m.Row(10, func() {
		m.Col(12, func() {
			textH3(m, "Rekap pekerjaan harian CCTV dan Altai", 0)
		})
	})

	// allListComputed = slice yang sudah di filter dan dimodifikasi isinya
	var allListComputed []dto.HistoryUnwindResponse
	var completeList []dto.HistoryUnwindResponse
	var progressList []dto.HistoryUnwindResponse
	var pendingList []dto.HistoryUnwindResponse

	// idTemp menyimpan id, karena akan banyak id yang sama, maka akan diambil history yang terakhir
	// urutan unwind dengan asumsi unwind sorted by updates.time 1 (pertama kali update tampil pertama)
	var idTemp string
	for _, history := range pdfVendorStruct.HistoryList {
		// skip jika waktu updatenya melebihi time end laporan
		if history.Updates.Time > pdfVendorStruct.End {
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

	startWita, _ = timegen.GetTimeWithYearWITA(pdfVendorStruct.Start)
	endWita, _ = timegen.GetTimeWithYearWITA(pdfVendorStruct.End)
	subtitle = fmt.Sprintf("Tanggal %s sd %s", startWita, endWita)

	if len(completeList) != 0 {
		buildTitleHeadingHistoryDailyView(m, " Pekerjaan Selesai", getTealColor())
		buildDailyHistoryVendorList(m, completeList)
	}

	if len(progressList) != 0 {
		buildTitleHeadingHistoryDailyView(m, " Pekerjaan Berjalan", getOrangeColor())
		buildDailyHistoryVendorList(m, progressList)
	}

	if len(pendingList) != 0 {
		buildTitleHeadingHistoryDailyView(m, " Pekerjaan Pending", getPinkColor())
		buildDailyHistoryVendorList(m, pendingList)
	}

	// simpan selesai ============================================
	err = m.OutputFileAndClose(fmt.Sprintf("static/pdf-vendor/%s.pdf", name))
	if err != nil {
		return err
	}
	return nil
}

func buildHeadingVendorDaily(m pdf.Maroto, subtitle string) error {
	var errTemp error
	m.Row(10, func() {

	})
	m.Row(20, func() {
		m.ColSpace(2)
		m.Col(8, func() {
			textH1(m, "Daily Maintenance Report")
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

func buildTitleHeadingView(m pdf.Maroto, title string, customColor color.Color) {
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
}

func buildTitleHeadingHistoryDailyView(m pdf.Maroto, title string, customColor color.Color) {
	m.SetBackgroundColor(customColor)
	m.Row(5, func() {
		m.Col(12, func() {
			m.Text(title, props.Text{
				Family:          consts.Courier,
				Style:           consts.Bold,
				Size:            9,
				Align:           consts.Left,
				VerticalPadding: 1,
				Color:           color.NewWhite(),
			})
		})
	})
	m.SetBackgroundColor(color.NewWhite())
}

func buildCCTVDailyView(m pdf.Maroto, cctv cctvDailyData, altai altaiDailyData) {

	// CCTV VIRTUAL ----- ALTAI VIRTUAL
	m.Row(15, func() {

		m.Col(6, func() {
			textH3(m, "CCTV", 3)
			textBodyItalic(m, fmt.Sprintf(" dicek pada %s", cctv.created), 8)
		})
		m.Col(6, func() {
			textH3(m, "ALTAI", 3)
			textBodyItalic(m, fmt.Sprintf(" dicek pada %s", altai.created), 8)
		})
	})

	// DATA 3-2-1     ------     DATA 3-2-1
	m.Row(30, func() {

		m.Col(3, func() {
			textBody(m, "- Total CCTV", 0)
			textBody(m, "- sudah di cek", 5)
			textBody(m, "- Kondisi ok", 10)
			textBody(m, "- Buram", 15)
			textBody(m, "- Offline", 20)
		})

		m.Col(2, func() {
			textBody(m, cctv.total, 0)
			textBody(m, cctv.checked, 5)
			textBody(m, cctv.ok, 10)
			textBody(m, cctv.blur, 15)
			textBody(m, cctv.offline, 20)
		})

		m.ColSpace(1)

		m.Col(3, func() {
			textBody(m, "- Total Altai", 0)
			textBody(m, "- Sudah di cek", 5)
			textBody(m, "- Kondisi Ok", 10)
			textBody(m, "- Offline", 15)
		})

		m.Col(2, func() {
			textBody(m, altai.total, 0)
			textBody(m, altai.checked, 5)
			textBody(m, altai.ok, 10)
			textBody(m, altai.offline, 15)
		})

		m.ColSpace(1)

	})

}

func buildVirtualTrouble(m pdf.Maroto, items []virtualTrouble) {
	m.SetBackgroundColor(getPastelColor())

	m.Row(5, func() {
		m.Col(12, func() {
			textBody(m, "Daftar perangkat bermasalah :", 0)
		})
	})

	tableHeading := []string{"Kategori", "Perangkat"}

	var contents [][]string
	for _, data := range items {
		contents = append(contents, []string{
			data.category,
			data.item,
		},
		)
	}
	m.TableList(tableHeading, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{4, 8},
		},
		ContentProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{4, 8},
			Family:    consts.Courier,
		},
		Align:              consts.Left,
		HeaderContentSpace: 1,
	})
	m.SetBackgroundColor(color.NewWhite())
}

func buildCCTVMonthlyView(m pdf.Maroto, cctv summaryMonthlyData, altai summaryMonthlyData) {

	// CCTV VIRTUAL ----- ALTAI VIRTUAL
	m.Row(15, func() {
		m.Col(6, func() {
			textH3(m, "CCTV", 3)
			textBodyItalic(m, fmt.Sprintf(" dimulai dari %s", cctv.created), 8)
		})
		m.Col(6, func() {
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

		m.ColSpace(1)

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

		m.ColSpace(1)

	})

}

func buildCCTVQuarterlyView(m pdf.Maroto, cctv summaryQuarterlyData, altai summaryQuarterlyData) {

	// CCTV VIRTUAL ----- ALTAI VIRTUAL
	m.Row(15, func() {
		m.Col(6, func() {
			textH3(m, "CCTV", 3)
			textBodyItalic(m, fmt.Sprintf(" dimulai dari %s", cctv.created), 8)
		})
		m.Col(6, func() {
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

		m.ColSpace(1)

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

		m.ColSpace(1)

	})

}

func buildCCTVQuarterlyViewNoAltai(m pdf.Maroto, cctv summaryQuarterlyData) {

	// CCTV VIRTUAL
	m.Row(15, func() {
		m.Col(6, func() {
			textH3(m, "CCTV", 3)
			textBodyItalic(m, fmt.Sprintf(" dimulai dari %s", cctv.created), 8)
		})
		m.ColSpace(6)
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

func buildDailyHistoryVendorList(m pdf.Maroto, dataList []dto.HistoryUnwindResponse) {
	tableHeading := []string{"Nama", "Keterangan", "Solusi", "Pengerjaan", "Update", "Oleh"}

	var contents [][]string
	for _, data := range dataList {
		updateAt, err := timegen.GetTimeWITA(data.Updates.Time)
		if err != nil {
			updateAt = "error"
		}

		contents = append(contents, []string{
			data.ParentName,
			data.Updates.Problem,
			data.Updates.ProblemResolve,
			sfunc.IntToTime(data.UpdatedAt, ""), // data UpdatedAt sudah diubah pada komputasi sebelumnya menjadi lama pengerjaan
			updateAt,
			strings.ToLower(data.UpdatedBy)},
		)
	}

	lightPurpleColor := getLightPurpleColor()
	m.TableList(tableHeading, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      7,
			GridSizes: []uint{2, 3, 3, 1, 1, 2},
		},
		ContentProp: props.TableListContent{
			Size:      6,
			GridSizes: []uint{2, 3, 3, 1, 1, 2},
		},
		Align:                consts.Left,
		AlternatedBackground: &lightPurpleColor,
		HeaderContentSpace:   1,
		Line:                 true,
	})
}
