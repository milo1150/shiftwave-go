package enum

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func (R *Role) ToString() string {
	switch *R {
	case RoleUser:
		return "user"
	case RoleAdmin:
		return "admin"
	default:
		return ""
	}
}

func ParseRole(text string) (*Role, bool) {
	switch text {
	case "user":
		role := RoleUser
		return &role, true
	case "admin":
		role := RoleAdmin
		return &role, true
	default:
		return nil, false
	}
}
