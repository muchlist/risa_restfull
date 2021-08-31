package stockpdf

import (
	"fmt"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/timegen"
	"strconv"
)

type PDFReq struct {
	Name      string
	StockList []dto.Stock
	Start     int64
	End       int64
}

func GenerateStockPDF(
	input PDFReq,
) error {

	m := pdf.NewMaroto(consts.Landscape, consts.A4)
	m.SetPageMargins(5, 10, 5)

	startWita, _ := timegen.GetTimeWithYearWITA(input.Start)
	endWita, _ := timegen.GetTimeWithYearWITA(input.End)
	subtitle := fmt.Sprintf("Tanggal %s sd %s", startWita, endWita)
	err := buildHeading(m, subtitle)
	if err != nil {
		return err
	}

	m.Row(5, func() {
		// space 5
	})
	if len(input.StockList) != 0 {
		buildTitleHeadingView(m, " Daftar Barang perlu restock", getPinkColor())
		buildStockList(m, input.StockList, input.Start, input.End)
		m.Row(5, func() {
			// space 5
		})
		m.Row(10, func() {
			m.Col(12, func() {
				textBody(m,
					fmt.Sprintf("NB : 1. Total penambahan dan pengurangan pada tabel adalah berdasarkan jarak waktu yang tertera pada sub-header dokumen, yaitu %s sampai dengan %s", startWita, endWita), 0)
				textBody(m, `        2. Kolom penambahan adalah seberapa banyak stok bertambah dari jumlah awalnya. Artinya bisa saja stok tersebut dikeluarkan lalu ditambahkan lagi karena tidak jadi dipakai.`, 4)
				textBody(m, `        3. Data diambil dari database aplikasi RISA yang mana setiap perubahan data akan memberikan notifikasi kepada semua IT di Regional Kalimantan.`, 8)
			})
		})
	} else {
		buildEmptyStock(m)
	}

	err = m.OutputFileAndClose(fmt.Sprintf("static/pdf-stock/%s.pdf", input.Name))
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
			textH1(m, "Stock Regional Kalimantan")
			textBodyCenter(m, subtitle, 12)
		})
		m.ColSpace(2)
	})
	return errTemp
}

func sumChangeStock(change []dto.StockChange, start, end int64) int {
	total := 0
	for _, c := range change {
		if c.Time >= start && c.Time <= end {
			total = total + c.Qty
		}
	}
	return total
}

func buildStockList(m pdf.Maroto, dataList []dto.Stock, start, end int64) {
	tableHeading := []string{"Nama Stok", "Kategori", "Sisa", "Penambahan atau Pengembalian", "Pengurangan", "Catatan"}
	var contents [][]string
	for _, data := range dataList {
		contents = append(contents, []string{
			data.Name,
			data.StockCategory,
			fmt.Sprintf("%d %s", data.Qty, data.Unit),
			strconv.Itoa(sumChangeStock(data.Increment, start, end)),
			strconv.Itoa(sumChangeStock(data.Decrement, start, end)),
			data.Note},
		)
	}

	lightPurpleColor := getLightPurpleColor()

	m.TableList(tableHeading, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{2, 1, 3, 2, 2, 2},
		},
		ContentProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{2, 1, 3, 2, 2, 2},
		},
		Align:                consts.Left,
		AlternatedBackground: &lightPurpleColor,
		HeaderContentSpace:   1,
		Line:                 true,
	})
}

func buildEmptyStock(m pdf.Maroto) {
	m.Row(10, func() {
		m.Col(12, func() {
			textBodyCenter(m, "Daftar stok perlu re-stok kosong", 0)
		})
	})
}
