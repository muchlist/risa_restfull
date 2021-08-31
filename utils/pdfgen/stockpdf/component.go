package stockpdf

import (
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func textH1(m pdf.Maroto, text string) {
	m.Text(text, props.Text{
		Top:         3,
		Style:       consts.Bold,
		Size:        18,
		Align:       consts.Center,
		Extrapolate: true,
		Color:       getDarkColor(),
	})

}

func textBody(m pdf.Maroto, text string, top float64) {
	m.Text(text, props.Text{
		Top:         top,
		Extrapolate: false,
		Size:        9,
		Color:       getDarkGreyColor(),
	})
}

func textBodyCenter(m pdf.Maroto, text string, top float64) {

	m.Text(text, props.Text{
		Top:         top,
		Extrapolate: false,
		Align:       consts.Center,
		Color:       getDarkGreyColor(),
	})

}

func getPinkColor() color.Color {
	return color.Color{
		Red:   255,
		Green: 51,
		Blue:  153,
	}
}

func getLightPurpleColor() color.Color {
	return color.Color{
		Red:   210,
		Green: 200,
		Blue:  230,
	}
}

func getDarkGreyColor() color.Color {
	return color.Color{
		Red:   83,
		Green: 83,
		Blue:  83,
	}
}

func getDarkColor() color.Color {
	return color.Color{
		Red:   36,
		Green: 36,
		Blue:  36,
	}
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
