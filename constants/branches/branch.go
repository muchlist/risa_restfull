package branches

const (
	Banjarmasin = "BANJARMASIN"
	Sampit      = "SAMPIT"
	Kumai       = "KUMAI"
	Kotabaru    = "KOTABARU"
	Batulicin   = "BATULICIN"
)

func GetBranchesAvailable() []string {
	return []string{Banjarmasin, Sampit, Kumai, Kotabaru, Batulicin}
}
