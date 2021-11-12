package ba

const (
	Equip     = "equip"
	Paragraph = "paragraph"
	Number    = "number"
	Bullet    = "bullet"
	Other     = "other"
)

func GetDescTypeAvailable() []string {
	return []string{Equip, Paragraph, Number, Bullet, Other}
}
