package location

import "github.com/muchlist/risa_restfull/constants/branches"

const (
	Regional  = "Regional"
	Trisakti  = "Trisakti"
	Tpkb      = "TPKB"
	Bmc       = "BMC"
	Marba     = "Marba"
	Pocc      = "POCC"
	Pulpis    = "Pulpis"
	Sampit    = "Sampit"
	Bagendang = "Bagendang"
	Kotabaru  = "Kotabaru"
	Batulicin = "Batulicin"
	Kumai     = "Kumai"
	Bumiharjo = "Bumiharjo"
	Lainnya   = "Lainnya"
)

func GetLocationAvailable() []string {
	return []string{
		Regional,
		Trisakti,
		Tpkb,
		Bmc,
		Marba,
		Pocc,
		Pulpis,
		Sampit,
		Bagendang,
		Kotabaru,
		Batulicin,
		Kumai,
		Bumiharjo,
		Lainnya,
	}
}

func GetLocationAvailableFrom(branch string) []string {
	switch branch {
	case branches.Banjarmasin:
		return []string{
			Regional,
			Trisakti,
			Tpkb,
			Bmc,
			Marba,
			Pocc,
			Pulpis,
			Lainnya,
		}
	case branches.Sampit:
		return []string{
			Sampit,
			Bagendang,
			Lainnya,
		}
	case branches.Kumai:
		return []string{
			Kumai,
			Bumiharjo,
			Lainnya,
		}
	case branches.Kotabaru:
		return []string{
			Kotabaru,
			Lainnya,
		}
	case branches.Batulicin:
		return []string{
			Batulicin,
			Lainnya,
		}
	default:
		return []string{}
	}
}
