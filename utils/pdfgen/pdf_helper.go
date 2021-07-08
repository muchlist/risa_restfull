package pdfgen

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

func textBodyCenter(m pdf.Maroto, text string, top float64) {

	m.Text(text, props.Text{
		Top:         top,
		Extrapolate: false,
		Align:       consts.Center,
		Color:       getDarkGreyColor(),
	})

}

func getTealColor() color.Color {
	return color.Color{
		Red:   3,
		Green: 166,
		Blue:  166,
	}
}

func getOrangeColor() color.Color {
	return color.Color{
		Red:   255,
		Green: 153,
		Blue:  51,
	}
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
