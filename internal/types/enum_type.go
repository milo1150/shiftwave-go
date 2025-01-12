package types

type Lang string

const (
	LangEN Lang = "EN"
	LangTH Lang = "TH"
	LangMY Lang = "MY"
)

func (l Lang) IsValid() bool {
	switch l {
	case LangEN, LangTH, LangMY:
		return true
	default:
		return false
	}
}
