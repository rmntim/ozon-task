package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type IDArray []uint

func (r IDArray) Value() (driver.Value, error) {
	if len(r) == 0 {
		return "[]", nil
	}

	sb := strings.Builder{}
	sb.WriteString("{")
	for _, id := range r {
		sb.WriteString(strconv.Itoa(int(id)))
		sb.WriteString(",")
	}
	sb.WriteString("}")
	return sb.String(), nil
}

func (r *IDArray) Scan(src interface{}) (err error) {
	var replies []uint
	switch src.(type) {
	case string:
		if src.(string)[1] == 'N' {
			*r = []uint{}
			return
		}
		err = json.Unmarshal([]byte(src.(string)), &replies)
	case []byte:
		if src.([]byte)[1] == 'N' {
			*r = []uint{}
			return
		}
		err = json.Unmarshal(src.([]uint8), &replies)
	default:
		err = errors.New("unknown type")
	}
	if err != nil {
		return
	}
	*r = replies
	return nil
}
