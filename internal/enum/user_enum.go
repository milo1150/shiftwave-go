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
