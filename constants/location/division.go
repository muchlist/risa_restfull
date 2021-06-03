package location

const (
	divSdm            = "SDM Umum"
	divQrm            = "QRM"
	divOperasional    = "Operasional"
	divKomersial      = "Komersial"
	divPelayananKapal = "Pelayanan Kapal"
	divTenik          = "Teknik"
	divKeuangan       = "Keuangan"
	divTIK            = "TIK"
	divTU             = "TU"
	divSmi            = "SMI"
	divGM             = "GM"
	divSEO            = "SEO"
	divServer         = "Server"
)

func GetDivisionAvailable() []string {
	return []string{
		divSdm,
		divQrm,
		divOperasional,
		divKomersial,
		divPelayananKapal,
		divTenik,
		divKeuangan,
		divTIK,
		divTU,
		divSmi,
		divGM,
		divSEO,
		divServer,
		Lainnya,
	}
}
