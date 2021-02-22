package roles

const (
	RoleAdmin   = "ADMIN"
	RoleIT      = "IT"
	RoleApprove = "APPROVE"
)

func GetRolesAvailable() []string {
	return []string{RoleAdmin, RoleIT, RoleApprove}
}
