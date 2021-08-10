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

type PDFVendorReq struct {
	Name            string
	HistoryList     dto.HistoryUnwindResponseList
	VendorCheckList []dto.VendorCheck
	Start           int64
	End             int64
}

func GeneratePDFVendor(
	pdfVendorStruct PDFVendorReq,
) error {
	// slice yang sudah di filter dan dimodifikasi isinya
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

	startWita, _ := timegen.GetTimeWithYearWITA(pdfVendorStruct.Start)
	endWita, _ := timegen.GetTimeWithYearWITA(pdfVendorStruct.End)
	subtitle := fmt.Sprintf("Tanggal %s sd %s", startWita, endWita)
	err := buildHeadingVendor(m, subtitle)
	if err != nil {
		return err
	}

	if len(completeList) != 0 {
		buildHistoryVendorList(m, completeList, " Completed", getTealColor())
	}

	if len(progressList) != 0 {
		buildHistoryVendorList(m, progressList, " Progress", getOrangeColor())
	}

	if len(pendingList) != 0 {
		buildHistoryVendorList(m, pendingList, " Pending", getPinkColor())
	}

	// Filtering check fisik
	var physicalCheckCCTVFiltered []dto.VendorCheckItemEmbed
	for _, checkParent := range pdfVendorStruct.VendorCheckList {
		if checkParent.IsVirtualCheck {
			continue
		}
		for _, check := range checkParent.VendorCheckItems {
			if !check.IsChecked {
				continue
			}
			if check.CheckedAt <= pdfVendorStruct.End && check.CheckedAt >= pdfVendorStruct.Start {
				physicalCheckCCTVFiltered = append(physicalCheckCCTVFiltered, check)
			}
		}
	}

	if len(physicalCheckCCTVFiltered) != 0 {
		m.AddPage()
		buildPhysicalCheckList(m, physicalCheckCCTVFiltered, " Cek Fisik")
	}

	err = m.OutputFileAndClose(fmt.Sprintf("static/pdf-vendor/%s.pdf", pdfVendorStruct.Name))
	if err != nil {
		return err
	}
	return nil
}

func buildHeadingVendor(m pdf.Maroto, subtitle string) error {
	var errTemp error
	m.Row(10, func() {

	})
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
			textH1(m, "Rekap Laporan Pekerjaan di Regional Kalimantan")
			textBodyCenter(m, subtitle, 12)
		})
		m.ColSpace(2)
	})
	return errTemp
}

func buildHistoryVendorList(m pdf.Maroto, dataList []dto.HistoryUnwindResponse, title string, customColor color.Color) {
	tableHeading := []string{"Nama", "Keterangan", "Solusi", "Status", "Pengerjaan", "Update", "Oleh"}

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
			enum.GetProgressString(data.Updates.CompleteStatus),
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
			GridSizes: []uint{2, 3, 3, 1, 1, 1, 1},
		},
		ContentProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{2, 3, 3, 1, 1, 1, 1},
		},
		Align:                consts.Left,
		AlternatedBackground: &lightPurpleColor,
		HeaderContentSpace:   1,
		Line:                 true,
	})
}

func buildPhysicalCheckList(m pdf.Maroto, dataList []dto.VendorCheckItemEmbed, title string) {
	tableHeading := []string{"No.", "CCTV", "Lokasi", "Offline", "Blur", "Pengecekan", "Oleh"}

	var contents [][]string
	for i, data := range dataList {
		updateAt, err := timegen.GetTimeWITA(data.CheckedAt)
		if err != nil {
			updateAt = "error"
		}

		blurText := ""
		offlineText := ""
		if data.IsBlur {
			blurText = "o"
		}
		if data.IsOffline {
			offlineText = "o"
		}

		contents = append(contents, []string{
			fmt.Sprintf("%03d\n", i+1),
			data.Name,
			data.Location,
			offlineText,
			blurText,
			updateAt,
			strings.ToLower(strings.Split(data.CheckedBy, " ")[0]),
		},
		)
	}

	lightPurpleColor := getLightPurpleColor()

	m.SetBackgroundColor(getTealColor())
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
			GridSizes: []uint{1, 3, 2, 1, 1, 2, 2},
		},
		ContentProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{1, 3, 2, 1, 1, 2, 2},
		},
		Align:                consts.Left,
		AlternatedBackground: &lightPurpleColor,
		HeaderContentSpace:   1,
		Line:                 true,
	})
}
