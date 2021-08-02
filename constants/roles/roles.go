package roles

const (
	RoleAdmin   = "ADMIN"
	RoleIT      = "IT"
	RoleApprove = "APPROVE"
	RoleVendor  = "VENDOR"
)

func GetRolesAvailable() []string {
	return []string{RoleAdmin, RoleIT, RoleApprove, RoleVendor}
}
