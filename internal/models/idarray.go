package models

import (
	"database/sql/driver"
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
		replies, err = parsePgArray(src.(string))
	case []byte:
		if src.([]byte)[1] == 'N' {
			*r = []uint{}
			return
		}
		replies, err = parsePgArray(string(src.([]byte)))
	default:
		err = errors.New("unknown type")
	}
	if err != nil {
		return
	}
	*r = replies
	return nil
}

func parsePgArray(array string) ([]uint, error) {
	var replies []uint
	nums := array[1 : len(array)-1]
	for _, num := range strings.Split(nums, ",") {
		id, err := strconv.Atoi(num)
		if err != nil {
			return nil, err
		}
		replies = append(replies, uint(id))
	}
	return replies, nil
}
