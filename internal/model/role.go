package model

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCompany  Role = "company"
	RoleCustomer Role = "customer"
)

func (r Role) String() string {
	return string(r)
}

func IsValidRole(r string) bool {
	switch Role(r) {
	case RoleAdmin, RoleCompany, RoleCustomer:
		return true
	}
	return false
}
