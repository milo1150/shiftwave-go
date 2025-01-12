package types

import (
	"errors"
)

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

func ParseLang(str string) (*Lang, error) {
	switch str {
	case "TH":
		l := LangTH
		return &l, nil
	case "EN":
		l := LangEN
		return &l, nil
	case "MY":
		l := LangMY
		return &l, nil
	default:
		return nil, errors.New("parse Lang error")
	}
}
