package models

import "strings"

type String string

func (s *String) UnmarshalJSON(bs []byte) (err error) {
	buf := string(bs)
	buf = strings.TrimSpace(buf)
	switch buf[0] {
	case '"':
		if len(buf[1 : len(buf)-1]) == 0 {
			return
		}
		buf = buf[1 : len(buf)-1]
		*s = String(buf)
	default:
		*s = String(buf)
	}
	return
}
